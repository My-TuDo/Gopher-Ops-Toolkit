package prober

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/semaphore"
)

// maxConcurrentProbes 控制单个 MultiRoundProbe 并发探测的最大 goroutine 数量
const maxConcurrentProbes = 5

// Prober 接口 — 所有探测器都实现这个接口
type Prober interface {
	Probe(target string, timeout time.Duration) Result
	Name() string
}

// Probers 全局注册表
var Probers = map[string]Prober{
	"tcp":  &TCPProber{},
	"http": &HTTPProber{},
	"dns":  &DNSProber{},
}

// ProbeStatus 探测状态枚举
type ProbeStatus int8

const (
	StatusUnknown ProbeStatus = iota
	StatusHealthy
	StatusUnhealthy
)

var statusStrings = map[ProbeStatus]string{
	StatusUnknown:   "未知",
	StatusHealthy:   "健康",
	StatusUnhealthy: "不健康",
}

func (s ProbeStatus) String() string {
	return statusStrings[s]
}

// StatusFromResult 从探测结果中提取状态，兼容中文老字段
func StatusFromResult(res Result) ProbeStatus {
	switch res.Status {
	case "健康":
		return StatusHealthy
	case "不健康":
		return StatusUnhealthy
	default:
		return StatusUnknown
	}
}

// MultiRoundProbe 对探针执行多轮并发探测，返回聚合结果。
//
// 设计要点：
//   - 使用 golang.org/x/sync/semaphore 加权信号量控制并发度，
//     避免大规模 goroutine 堆积在 channel send 上。
//   - 熔断阈值 = maxConcurrentProbes / 2，失败数达到立即 Cancel 剩余探测。
//   - 所有 goroutine 通过 context 联动取消，不存在 goroutine 泄漏。
//   - 结果通过带缓冲 channel 收集，避免发送者在 wg.Wait 完成后阻塞。
func MultiRoundProbe(p Prober, target string, timeout time.Duration, rounds int) Result {
	if rounds <= 1 {
		return p.Probe(target, timeout)
	}

	// ——— 熔断阈值 ———
	failureThreshold := int32(maxConcurrentProbes / 2)
	if failureThreshold < 1 {
		failureThreshold = 1
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 加权信号量 — 优雅控制并发度，不会让 goroutine 卡在 channel send
	sem := semaphore.NewWeighted(int64(maxConcurrentProbes))
	var wg sync.WaitGroup
	var (
		failureCount atomic.Int32
		successCount atomic.Int32
		totalLatency atomic.Int64
	)

	for i := 0; i < rounds; i++ {
		// 熔断检查：失败数已达阈值，不再 launch 新的 goroutine
		if failureCount.Load() >= failureThreshold {
			cancel()
			break
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			// 获取信号量令牌（支持 context 取消）
			if err := sem.Acquire(ctx, 1); err != nil {
				// context 已取消，安静退出
				return
			}
			defer sem.Release(1)

			res := p.Probe(target, timeout)

			status := StatusFromResult(res)

			if status == StatusHealthy {
				successCount.Add(1)
				totalLatency.Add(res.Latency)
			} else {
				// 熔断：失败数达到阈值则取消所有剩余探测
				if failureCount.Add(1) >= failureThreshold {
					cancel()
				}
			}
		}()
	}

	wg.Wait()

	name := p.Name()
	ok := successCount.Load()
	if ok == 0 {
		return Result{
			Name:    name,
			Target:  target,
			Status:  StatusUnhealthy.String(),
			Error:   fmt.Sprintf("%d 轮探测全部失败", rounds),
			Latency: 0,
		}
	}

	avgLatency := totalLatency.Load() / int64(ok)
	return Result{
		Name:    name,
		Target:  target,
		Status:  StatusHealthy.String(),
		Detail:  fmt.Sprintf("%d/%d 轮探测成功，平均延迟: %d ms", ok, rounds, avgLatency),
		Latency: avgLatency,
	}
}
