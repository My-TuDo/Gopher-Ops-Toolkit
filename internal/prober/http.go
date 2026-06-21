package prober

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type HTTPProber struct {
	Method  string            // HTTP 方法，如 GET、POST 等
	Headers map[string]string // 可选的 HTTP 头部信息
}

func (h HTTPProber) Probe(target string, timeout time.Duration) Result {
	// 计算开始时间
	start := time.Now()

	// 先进行数据清洗
	// 防止用户输入带有多余空格的地址导致探测失败
	cleanedTarget := strings.TrimSpace(target)

	// 判断是否以 "http://" 或 "https://" 开头，如果没有则默认加上 "http://"
	if !strings.HasPrefix(cleanedTarget, "http://") && !strings.HasPrefix(cleanedTarget, "https://") {
		// 默认使用 http 协议
		cleanedTarget = "http://" + cleanedTarget
	}

	// 创建 HTTP 客户端
	client := http.Client{}

	// 创建动态超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// h.Method 如果没有设置，默认使用 GET 方法
	method := h.Method
	if method == "" {
		method = "GET"
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, method, cleanedTarget, nil)
	if err != nil {
		return Result{
			Name:    "HTTP探测",
			Target:  cleanedTarget,
			Status:  "不健康",
			Error:   err.Error(),
			Latency: time.Since(start).Milliseconds(),
		}
	}

	// 设置请求头
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		return Result{
			Name:    "HTTP探测",
			Target:  cleanedTarget,
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
			Target:  cleanedTarget,
			Status:  "健康",
			Detail:  "HTTP状态码正常",
			Latency: time.Since(start).Milliseconds(),
		}
	}

	return Result{
		Name:    "HTTP探测",
		Target:  cleanedTarget,
		Status:  "不健康",
		Error:   "HTTP状态码异常",
		Latency: time.Since(start).Milliseconds(),
	}
}
