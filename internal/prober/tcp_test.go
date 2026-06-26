package prober

import (
	"net"
	"strings"
	"testing"
	"time"
)

// 测试 TCP 探测器的基本功能
func TestTCPProber(t *testing.T) {
	// 拉起一个 TCP 服务器用于测试
	ln, err := net.Listen("tcp", "127.0.0.1:0") // 监听随机端口
	if err != nil {
		t.Fatalf("无法启动测试 TCP 服务器: %v", err)
	}
	defer ln.Close()

	// 获取监听的端口
	addr := ln.Addr().String()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return // 监听关闭时自动退出
			}
			// 握手成功后，返回一字节数据或在直接挂断，模拟正常的网络连接
			_ = conn.SetDeadline(time.Now().Add(1 * time.Second)) // 设置超时
			conn.Close()
		}
	}()

	// 创建测试矩阵
	tests := []struct {
		name         string
		target       string
		timeout      time.Duration
		expectStatus string
		expectErrCtx string
	}{
		{
			name:         "端口开放，模拟正常 TCP 握手",
			target:       addr,
			timeout:      3 * time.Second,
			expectStatus: "健康",
		},
		{
			name:         "端口未开放，模拟拒绝连接",
			target:       "127.0.0.1:59999", // 给一个不可能开放的端口
			timeout:      3 * time.Second,
			expectStatus: "不健康",
			expectErrCtx: "refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := TCPProber{}
			res := p.Probe(tt.target, tt.timeout)

			// 验证状态
			if res.Status != tt.expectStatus {
				t.Errorf("期望状态 %s，得到 %s", tt.expectStatus, res.Status)
			}

			// 验证错误信息
			if tt.expectErrCtx != "" && !strings.Contains(res.Error, tt.expectErrCtx) {
				t.Errorf("期望错误信息包含 %s，得到 %s", tt.expectErrCtx, res.Error)
			}
		})
	}
}
