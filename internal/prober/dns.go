package prober

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type DNSProber struct {
	RecordType string // 可选参数，指定 DNS 记录类型，如 A、AAAA、CNAME 等
	DNSServer  string // 可选参数，指定 DNS 服务器地址，如 8.8.8.8
}

func (d DNSProber) Name() string {
	return "DNS探测"
}

func (d DNSProber) Probe(target string, timeout time.Duration) Result {
	// 开始计时
	start := time.Now()

	// 设置动态超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 初始化一个原生的 GO 解析器
	resolver := &net.Resolver{PreferGo: true}
	domain := target // 默认为全域名

	// 如果指定了 DNS 服务器，替换默认拨号
	if d.DNSServer != "" {
		dnsServer := d.DNSServer

		if !strings.Contains(dnsServer, ":") {
			dnsServer = dnsServer + ":53" // 默认 DNS 服务器端口
		}

		// 把 go 默认的拨号行为，改成连指定的 DNS 服务器
		resolver.Dial = func(ctx context.Context, network, address string) (net.Conn, error) {
			dialer := net.Dialer{Timeout: timeout}
			return dialer.DialContext(ctx, "udp", dnsServer)
		}
	}

	var detail string
	var err error

	switch d.RecordType {
	case "A", "AAAA", "":
		var ips []net.IP
		switch d.RecordType {
		case "A":
			ips, err = resolver.LookupIP(ctx, "ip4", domain)

		case "AAAA":
			ips, err = resolver.LookupIP(ctx, "ip6", domain)

		default:
			ips, err = resolver.LookupIP(ctx, "ip", domain)
		}
		if err == nil {
			detail = fmt.Sprintf("Resolved to: %s", joinIPs(ips))
		}

	case "CNAME":
		var cname string
		cname, err = resolver.LookupCNAME(ctx, domain)
		if err == nil {
			detail = fmt.Sprintf("CNAME: %s", cname)
		}

	case "MX":
		var mxRecords []*net.MX
		mxRecords, err = resolver.LookupMX(ctx, domain)
		if err == nil {
			var sb []string
			for _, mx := range mxRecords {
				sb = append(sb, fmt.Sprintf("%s (优先%d)", mx.Host, mx.Pref))
			}
			detail = fmt.Sprintf("MX: %s", strings.Join(sb, ", "))
		}

	case "NS":
		var nsRecords []*net.NS
		nsRecords, err = resolver.LookupNS(ctx, domain)
		if err == nil {
			var sb []string
			for _, ns := range nsRecords {
				sb = append(sb, ns.Host)
			}
			detail = fmt.Sprintf("NS: %s", strings.Join(sb, ", "))
		}

	case "TXT":
		var txtRecords []string
		txtRecords, err = resolver.LookupTXT(ctx, domain)
		if err == nil {
			detail = fmt.Sprintf("TXT:%s", strings.Join(txtRecords, ", "))
		}

	default:
		return Result{
			Name:    d.Name(),
			Target:  target,
			Status:  "不健康",
			Error:   fmt.Sprintf("不支持的 DNS 记录类型: %s", d.RecordType),
			Latency: time.Since(start).Milliseconds(),
		}
	}

	elapsed := time.Since(start).Milliseconds()
	if err != nil {
		return Result{
			Name:    d.Name(),
			Target:  target,
			Status:  "不健康",
			Error:   err.Error(),
			Latency: elapsed,
		}
	}
	return Result{
		Name:    d.Name(),
		Target:  target,
		Status:  "健康",
		Latency: elapsed,
		Detail:  detail,
	}
}

// 辅助函数：将 []net.IP 转为字符串
func joinIPs(ips []net.IP) string {
	strs := make([]string, len(ips))
	for i, ip := range ips {
		strs[i] = ip.String()
	}
	return strings.Join(strs, ", ")
}
