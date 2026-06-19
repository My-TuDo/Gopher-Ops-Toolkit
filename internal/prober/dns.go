package prober

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type DNSProber struct{}

func (d DNSProber) Probe(target string, timeout time.Duration) Result {
	// 开始计时
	start := time.Now()

	// 设置动态超时控制
	ctx, cancel := time.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 初始化一个原生的 GO 解析器
	resolver := &net.Resolver{PreferGo: true}
	domain := target // 默认为全域名

	// 如果有 @ 符号，说明用户指定了 DNS 服务器
	if strings.Contains(target, "@") {
		parts := strings.Split(target, "@")
		domain = parts[0]     // 要查的 web 域名
		dnsServer := parts[1] // 外部的 DNS 服务器地址

		// 把 go 默认的拨号行为，改成连指定的 DNS 服务器
		resolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{Timeout: timeout}
			return dialer.DialContext(ctx, "udp", dnsServer)
		}
	}

	// 执行 DNS 解析
	ips, err := resolver.LookupHost(ctx, domain)
	elapsed := time.Since(start).Milliseconds()

	if err != nil {
		return Result{
			Name:    "DNS探测",
			Target:  target,
			Status:  "不健康",
			Error:   err.Error(),
			Latency: elapsed,
		}
	}

	return Result{
		Name:    "DNS探测",
		Target:  target,
		Status:  "健康",
		Latency: elapsed,
		Detail:  fmt.Sprintf("Resolved to: %s", strings.Join(ips, ", ")), // 解析到的 IP 地址列表作为 Detail 返回
	}
}
