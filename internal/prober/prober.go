package prober

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// maxConcurrentProbes 控制并发探测的最大 goroutine 数量
	maxConcurrentProbes = 5
)

// 定义一个接口，所有探测器都实现这个接口
type Prober interface {
	Probe(target string, timeout time.Duration) Result
	Name() string
}

// 全局注册表
var Probers = map[string]Prober{
	"tcp":  &TCPProber{},
	"http": &HTTPProber{},
	"dns":  &DNSProber{},
}

var globalSem = make(chan struct{}, 20) // 先硬编码，全局最多 20 个并发 goroutine

// MultiRoundProbe 对探针执行多轮并发探测，返回平均值。
//
// 特性：
//   - 信号量限流：最多 maxConcurrentProbes 个 goroutine 同时运行
//   - 熔断机制：失败次数超过阈值时自动取消剩余探测
//   - 优雅取消：通过 context 通知所有 goroutine 退出
func MultiRoundProbe(p Prober, target string, timeout time.Duration, rounds int) Result {
	globalSem <- struct{}{}        // 获取全局令牌
	defer func() { <-globalSem }() // 释放全局令牌

	if rounds <= 1 {
		return p.Probe(target, timeout)
	}

	// 熔断阈值：失败超过并发上限的一半即触发熔断
	failureThreshold := int32(maxConcurrentProbes / 2)
	if failureThreshold < 1 {
		failureThreshold = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type roundResult struct {
		latency int64
		success bool
	}

	resultCh := make(chan roundResult, rounds)
	sem := make(chan struct{}, maxConcurrentProbes) // 信号量：控制并发数
	var wg sync.WaitGroup
	var failureCount atomic.Int32
	var cancelled atomic.Bool // 标记是否已取消

	for i := 0; i < rounds; i++ {
		// 令牌获取前检查熔断
		if failureCount.Load() >= failureThreshold {
			cancel()
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			// 获取信号量令牌（支持上下文取消）
			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				return
			}
			defer func() { <-sem }()

			// 执行探测前再次检查上下文
			if ctx.Err() != nil {
				return
			}

			res := p.Probe(target, timeout)

			if cancelled.Load() {
				return
			}
			success := res.Status == "健康"

			if !success {
				// 失败计数达到阈值时触发熔断
				if failureCount.Add(1) >= failureThreshold {
					cancel()
					cancelled.Store(true)
				}
			}

			resultCh <- roundResult{
				latency: res.Latency,
				success: success,
			}
		}()
	}

	wg.Wait()
	close(resultCh)

	var totalLatency int64
	successCount := 0
	for r := range resultCh {
		if r.success {
			totalLatency += r.latency
			successCount++
		}
	}

	name := p.Name()
	if successCount == 0 {
		return Result{
			Name:    name,
			Target:  target,
			Status:  "不健康",
			Error:   fmt.Sprintf("%d 轮探测全部失败", rounds),
			Latency: 0,
		}
	}

	return Result{
		Name:    name,
		Target:  target,
		Status:  "健康",
		Detail:  fmt.Sprintf("%d/%d 轮探测成功，平均延迟: %d ms", successCount, rounds, totalLatency/int64(successCount)),
		Latency: totalLatency / int64(successCount),
	}
}
