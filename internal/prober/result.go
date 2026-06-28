package prober

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	Name    string `json:"name"`              // 项目名称
	Target  string `json:"target"`            // 目标地址
	Status  string `json:"status"`            // 健康状态
	Error   string `json:"error,omitempty"`   // 错误信息（如果有）
	Detail  string `json:"detail,omitempty"`  // 详细信息（可选）
	Latency int64  `json:"latency,omitempty"` // 响应时间（如果健康）
}

// String 格式化输出结果
func (r Result) String() string {
	if r.Error != "" {
		return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s, 错误: %s", r.Name, r.Target, r.Status, r.Error)
	}
	return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s, 详情：%s", r.Name, r.Target, r.Status, r.Detail)
}

const (
	// maxConcurrentProbes 控制并发探测的最大 goroutine 数量
	maxConcurrentProbes = 5
)

// MultiRoundProbe 对探针执行多轮并发探测，返回平均值。
//
// 特性：
//   - 信号量限流：最多 maxConcurrentProbes 个 goroutine 同时运行
//   - 熔断机制：失败次数超过阈值时自动取消剩余探测
//   - 优雅取消：通过 context 通知所有 goroutine 退出
func MultiRoundProbe(p Prober, target string, timeout time.Duration, rounds int) Result {
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
			success := res.Status == "健康"

			if !success {
				// 失败计数达到阈值时触发熔断
				if failureCount.Add(1) >= failureThreshold {
					cancel()
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
