package prober

import (
	"context"
	"net/http"
	"time"
)

type HTTPProber struct{}

func (h HTTPProber) Probe(target string, timeout time.Duration) Result {
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

}
