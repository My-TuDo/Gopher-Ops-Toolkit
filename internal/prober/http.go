package prober

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPProber struct {
	Method  string            // HTTP 方法，如 GET、POST 等
	Headers map[string]string // 可选的 HTTP 头部信息
	Keyword string            // 期望 body 中包含的关键字
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
	bash := Result{
		Name:   "HTTP探测",
		Target: target,
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
		bash.Status = "不健康"
		bash.Error = fmt.Sprintf("创建 HTTP 请求失败: %v", err)
		bash.Latency = time.Since(start).Milliseconds()
		return bash
	}

	// 设置请求头
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}

	// 发送 HTTP 请求
	resp, err := client.Do(req)
	if err != nil {
		bash.Status = "不健康"
		bash.Error = fmt.Sprintf("发送 HTTP 请求失败: %v", err)
		bash.Latency = time.Since(start).Milliseconds()
		return bash
	}
	defer resp.Body.Close()

	// 判断 HTTP 状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bash.Detail = fmt.Sprintf("HTTP 状态码: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		// 如果有关键字检测需求，读 body 检查
		if h.Keyword != "" {
			body, _ := io.ReadAll(resp.Body)
			if !strings.Contains(string(body), h.Keyword) {
				bash.Status = "不健康"
				bash.Error = fmt.Sprintf("HTTP 响应体不包含关键字: %s", h.Keyword)
				bash.Latency = time.Since(start).Milliseconds()
				return bash
			}
			bash.Detail += fmt.Sprintf("| 响应体大小：%d 字节", len(body))
		}
		bash.Status = "健康"
		bash.Latency = time.Since(start).Milliseconds()
		return bash
	}

	bash.Status = "不健康"
	bash.Error = fmt.Sprintf("HTTP 状态码异常: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	bash.Latency = time.Since(start).Milliseconds()
	return bash
}
