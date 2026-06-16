package prober

import (
	"net"
	"time"
)

type DNSProber struct{}

func (d DNSProber) Probe(target string, timeout time.Duration) Result {
	// 开始计时
	start := time.Now()

	// 先实现对本机 DNS 服务器的探测检查
	ips, err := net.LookupIP(target)

	// 计算响应时间
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return Result{
			Name:    "DNS探测",
			Target:  target,
			Status:  "不健康",
			Error:   err.Error(),
			Latency: latency,
		}
	}

	// 解析成功但没有 IP 地址
	if len(ips) == 0 {
		return Result{
			Name:    "DNS探测",
			Target:  target,
			Status:  "不健康",
			Error:   "没有解析到 IP 地址",
			Latency: latency,
		}
	}

	// 探测成功
	return Result{
		Name:    "DNS探测",
		Target:  target,
		Status:  "健康",
		Detail:  "解析成功: " + ips[0].String(),
		Latency: latency,
	}
}
