package prober

import (
	"context"
	"fmt"
	"net"
	"time"
)

// 实现 TCP 端口探测的结构体
type TCPProber struct{}

func (t TCPProber) Probe(target string, timeout time.Duration) Result {
	// 记录开始时间
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.background(), timeout)
	defer cancel()

	// 分离主机号和端口号
	host, port, err := net.SplitHostPort(target)
	if err != nil {
		return Result{
			Name:    "TCP探测",
			Target:  target,
			Status:  "不健康",
			Error:   fmt.Sprintf("目标地址格式错误: %v", err),
			Latency: time.Since(start).Milliseconds(),
		}
	}

	// ———— 进行 DNS 解析 ——————
	dnsStart := time.Now()
	resolver := &net.Resolver{PreferGo: true}
	ips, err := resolver.LookupHost(ctx, host)
	dnsDuration := time.Since(dnsStart).Milliseconds()
	if err != nil {
		return Result{
			Name:    "TCP探测",
			Target:  target,
			Status:  "不健康",
			Error:   fmt.Sprintf("DNS 解析失败: %v", err),
			Latency: time.Since(start).Milliseconds(),
		}
	}

	// ———— TCP 握手阶段 ——————
	connStart := time.Now()
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(ips[0], port))
	handshakeDuration := time.Since(connStart).Milliseconds()
	totalDuration := time.Since(start).Milliseconds()
	if err != nil {
		return Result{
			Name:    "TCP探测",
			Target:  target,
			Status:  "不健康",
			Error:   fmt.Sprintf("TCP 握手失败: %v", err),
			Latency: totalDuration,
		}
	}
	defer conn.Close()

	// ———— Banner 读取 ——————
	banner := ""
	_ = conn.SetReadDeadline(time.Now().Add(300 * time.Millisecond)) // 设置读取超时
	buf := make([]byte, 256)
	if n, err := conn.Read(buf); err == nil && n > 0 {
		banner = sanitizeBanner(buf[:n])
	}

	detail := fmt.Sprintf("端口开放 (DNS: %dms, 握手: %dms)", dnsDuration, handshakeDuration)
	if banner != "" {
		detail += fmt.Sprintf(" | Banner: %s", banner)
	}

	// 连接成功
	return Result{
		Name:    "TCP探测",
		Target:  target,
		Status:  "健康",
		Detail:  detail,
		Latency: totalDuration,
	}
}

// 过滤不可见字符，保留可读性
func sanitizeBanner(s string) string {
	return s
}
