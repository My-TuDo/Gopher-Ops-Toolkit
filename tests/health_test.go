package tests

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// 编译二进制
func TestMain(m *testing.M) {
	if err := exec.Command("go", "build", "-o", "../ops-toolkit", "..").Run(); err != nil {
		os.Stderr.WriteString("编译失败: " + err.Error() + "\n")
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func TestHealthCheck(t *testing.T) {
	// 启动本地服务器

	// 启动TCP服务器
	ln, err := net.Listen("tcp", "127.0.0.1:0") // 监听随机端口
	if err != nil {
		t.Fatalf("无法启动测试 TCP 服务器: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return // 监听关闭时自动退出
			}
			conn.Close() // 立即关闭连接，模拟服务响应
		}
	}()
	// 获取监听的端口
	port := ln.Addr().(*net.TCPAddr).Port

	// 编译并执行子进程
	cmd := exec.Command("../ops-toolkit", "health",
		"--config", "../configs/config.yaml",
		"--type", "tcp",
		"--host", "127.0.0.1",
		"--port", fmt.Sprintf("%d", port),
	)
	output, err := cmd.CombinedOutput()

	// 断言
	if err != nil {
		t.Fatalf("子进程执行失败(预期成功): %v\n输出: %s", err, string(output))
	}
	if !strings.Contains(string(output), "健康") {
		t.Errorf("期望输出包含 '健康'，但实际输出为: %s", string(output))
	}

}

func TestHealth_TCP_Refused(t *testing.T) {
	cmd := exec.Command("../ops-toolkit", "health",
		"--config", "../configs/config.yaml",
		"--type", "tcp",
		"--host", "127.0.0.1",
		"--port", "59999",
	)
	output, err := cmd.CombinedOutput()

	// 断言
	if err == nil {
		t.Fatalf("期望子进程非零退出，但实际成功")
	}
	if !strings.Contains(string(output), "不健康") {
		t.Errorf("期望输出包含 '不健康'，但实际输出为: %s", string(output))
	}
}
