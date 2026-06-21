package prober

import (
	"fmt"
	"sync"
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
	return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s", r.Name, r.Target, r.Status)
}

// MultiRoundProbe 对探针执行多轮并发探测，返回平均值
func MultiRoundProbe(p Prober, target string, timeout time.Duration, rounds int) Result {
	if rounds <= 1 {
		return p.Probe(target, timeout)
	}

	type roundResult struct {
		latency int64
		success bool
	}

	reusltCh := make(chan roundResult, rounds)
	var wg sync.WaitGroup
	// 每轮分配 timeout
	for i := 0; i < rounds; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res := p.Probe(target, timeout)
			reusltCh <- roundResult{
				latency: res.Latency,
				success: res.Status == "健康",
			}
		}()
	}

	wg.Wait()
	close(reusltCh)

	var totalLatency int64
	successCount := 0
	for r := range reusltCh {
		if r.success {
			totalLatency += r.latency
			successCount++
		}
	}

	firstRes := p.Probe(target, timeout)
	if successCount == 0 {
		return Result{
			Name:    firstRes.Name,
			Target:  target,
			Status:  "不健康",
			Error:   fmt.Sprintf("%d 轮探测全部失败", rounds),
			Latency: 0,
		}
	}

	return Result{
		Name:    firstRes.Name,
		Target:  target,
		Status:  "健康",
		Detail:  fmt.Sprintf("%d 轮探测成功，平均延迟: %d ms", successCount, totalLatency/int64(successCount)),
		Latency: totalLatency / int64(successCount),
	}
}

// 定义一个接口，所有探测器都实现这个接口
type Prober interface {
	Probe(target string, timeout time.Duration) Result
}

// 全局注册表
var Probers = map[string]Prober{
	"tcp":  &TCPProber{},
	"http": &HTTPProber{},
	"dns":  &DNSProber{},
}
