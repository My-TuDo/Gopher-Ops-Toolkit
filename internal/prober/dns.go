package prober

import "time"

type DNSProber struct{}

func (d DNSProber) Probe(target string, timeout time.Duration) Result {

	return Result{
		Name:    "DNS探测",
		Target:  target,
		Status:  "不健康",
		Error:   "DNS探测功能尚未实现",
		Latency: 0,
	}
}
