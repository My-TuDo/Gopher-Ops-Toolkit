package prober

import (
	"net"
	"time"
)

// 实现 TCP 端口探测的结构体
type TCPProber struct{}

func (t TCPProber) Probe(target string, timeout time.Duration) Result {
	// 记录开始时间
	start := time.Now()

	// 尝试连接 TCP 端口
	conn, err := net.DialTimeout("tcp", target, timeout)

	// 计算响应时间
	latency := time.Since(start).Milliseconds()

	// 连接失败
	if err != nil {
		return Result{
			Name:    "TCP探测",
			Target:  target,
			Status:  "不健康",
			Error:   err.Error(),
			Latency: latency,
		}
	}
	defer conn.Close()

	// 连接成功
	return Result{
		Name:    "TCP探测",
		Target:  target,
		Status:  "健康",
		Detail:  "端口开放",
		Latency: latency,
	}
}
