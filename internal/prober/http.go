package prober

import (
	"context"
	"net/http"
	"time"
)

type HTTPProber struct{}

func (h HTTPProber) Probe(target string, timeout time.Duration) Result {
	target = "http://" + target

	// 计算开始时间
	start := time.Now()

	// 创建 HTTP 客户端
	client := http.Client{}

	// 创建动态超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", target, nil)
	if err != nil {
		return Result{
			Name:    "HTTP探测",
			Target:  target,
			Status:  "不健康",
			Error:   err.Error(),
			Latency: time.Since(start).Milliseconds(),
		}
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		return Result{
			Name:    "HTTP探测",
			Target:  target,
			Status:  "不健康",
			Error:   err.Error(),
			Latency: time.Since(start).Milliseconds(),
		}
	}
	defer resp.Body.Close()

	// 判断 HTTP 状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return Result{
			Name:    "HTTP探测",
			Target:  target,
			Status:  "健康",
			Detail:  "HTTP状态码正常",
			Latency: time.Since(start).Milliseconds(),
		}
	}

	return Result{
		Name:    "HTTP探测",
		Target:  target,
		Status:  "不健康",
		Error:   "HTTP状态码异常",
		Latency: time.Since(start).Milliseconds(),
	}
}
