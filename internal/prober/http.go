package prober

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	// maxBodySize 限制读取 HTTP 响应体的最大字节数（10MB）
	maxBodySize = 10 * 1024 * 1024
)

// defaultHTTPClient 是包级复用的 HTTP 客户端，带连接池配置。
// 每次 Probe 不再新建 http.Client，避免 Transport 连接池无法复用。
var defaultHTTPClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConns:        20,
		MaxIdleConnsPerHost: 5,
		IdleConnTimeout:     90 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2: true,
	},
	// 不设置 Timeout，由 request context 控制超时
	Timeout: 0,
}

type HTTPProber struct {
	Method  string            // HTTP 方法，如 GET、POST 等
	Headers map[string]string // 可选的 HTTP 头部信息
	Keyword string            // 期望 body 中包含的关键字
}

func (h HTTPProber) Name() string {
	return "HTTP探测"
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
	base := Result{
		Name:   h.Name(),
		Target: target,
	}
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
		base.Status = "不健康"
		base.Error = fmt.Sprintf("创建 HTTP 请求失败: %v", err)
		base.Latency = time.Since(start).Milliseconds()
		return base
	}

	// 设置请求头
	for key, value := range h.Headers {
		req.Header.Set(key, value)
	}

	// 使用复用客户端发送 HTTP 请求
	resp, err := defaultHTTPClient.Do(req)
	if err != nil {
		base.Status = "不健康"
		base.ErrorCode = ErrCodeSystemError
		base.Error = fmt.Sprintf("发送 HTTP 请求失败: %v", err)
		base.Latency = time.Since(start).Milliseconds()
		return base
	}

	// 确保连接归还连接池：先耗尽 body，再 close
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
	}()

	// 判断 HTTP 状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		base.Detail = fmt.Sprintf("HTTP 状态码: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
		// 如果有关键字检测需求，读 body 检查（限制最大读取量）
		if h.Keyword != "" {
			limitedReader := io.LimitReader(resp.Body, maxBodySize)
			body, readErr := io.ReadAll(limitedReader)
			if readErr != nil {
				base.Status = "不健康"
				base.ErrorCode = ErrCodeSystemError
				base.Error = fmt.Sprintf("读取响应体失败: %v", readErr)
				base.Latency = time.Since(start).Milliseconds()
				return base
			}
			if !strings.Contains(string(body), h.Keyword) {
				base.Status = "不健康"
				base.ErrorCode = ErrCodeSystemError
				base.Error = fmt.Sprintf("HTTP 响应体不包含关键字: %s", h.Keyword)
				base.Latency = time.Since(start).Milliseconds()
				return base
			}
			base.Detail += fmt.Sprintf(" | 响应体大小：%d 字节", len(body))
		}
		base.Status = "健康"
		base.Latency = time.Since(start).Milliseconds()
		return base
	}

	base.Status = "不健康"
	base.ErrorCode = ErrCodeSystemError
	base.Error = fmt.Sprintf("HTTP 状态码异常: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	base.Latency = time.Since(start).Milliseconds()
	return base
}
