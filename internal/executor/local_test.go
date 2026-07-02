package executor

import (
	"context"
	"strings"
	"testing"
	"time"
)

// 测试基本功能：正常执行并获取输出
func TestLocalExecutor_Basic(t *testing.T) {
	exec := NewLocalExecutor()
	ctx := context.Background()

	resultCh, err := exec.Exec(ctx, "echo Hello World")
	if err != nil {
		t.Fatalf("启动命令失败: %v", err)
	}

	var lines []string
	for res := range resultCh {
		if res.Done {
			if res.Err != nil {
				t.Fatalf("命令执行异常: %v", res.Err)
			}
			if res.ExitCode != 0 {
				t.Fatalf("期望退出码 0，得到 %d", res.ExitCode)
			}
		} else {
			lines = append(lines, res.Line)
		}
	}

	// 验证输出了 Hello World
	if len(lines) == 0 || lines[0] != "Hello World" {
		t.Fatalf("期望输出 'Hello World'，得到 %v", lines)
	}
}

// 测试超时：确认命令被杀死，且退出码非零
func TestLocalExecutor_Timeout(t *testing.T) {
	exec := NewLocalExecutor()
	// 超短超时，确保 sleep 1 被杀死
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	resultCh, err := exec.Exec(ctx, "sleep 1")
	if err != nil {
		t.Fatalf("启动命令失败: %v", err)
	}

	var doneRes *ExecResult
	for res := range resultCh {
		if res.Done {
			r := res
			doneRes = &r
		}
	}

	if doneRes == nil {
		t.Fatal("未收到 Done 信号")
	}

	// 超时杀掉进程，err 应该不为 nil
	if doneRes.Err == nil {
		t.Fatal("期望命令被超时杀死（err != nil），但命令正常结束了")
	}
	// 退出码应该非零（被信号杀死）
	if doneRes.ExitCode == 0 {
		t.Fatal("期望退出码非 0（被信号杀死），得到 0")
	}

	// 验证错误信息包含 "killed"（语义检查）
	if !strings.Contains(doneRes.Err.Error(), "killed") &&
		!strings.Contains(doneRes.Err.Error(), "signal") {
		t.Fatalf("期望错误包含 'killed' 或 'signal'，得到: %v", doneRes.Err)
	}
}

// 测试 stderr 分离：stdout 和 stderr 应该独立输出
func TestLocalExecutor_Stderr(t *testing.T) {
	exec := NewLocalExecutor()
	ctx := context.Background()

	// 同时输出到 stdout 和 stderr
	resultCh, err := exec.Exec(ctx, "echo 'stdout msg'; echo 'stderr msg' >&2")
	if err != nil {
		t.Fatalf("启动命令失败: %v", err)
	}

	var stdoutLines, stderrLines []string
	for res := range resultCh {
		if res.Done {
			continue
		}
		if res.IsStderr {
			stderrLines = append(stderrLines, res.Line)
		} else {
			stdoutLines = append(stdoutLines, res.Line)
		}
	}

	// 验证 stdout 内容
	if len(stdoutLines) != 1 || stdoutLines[0] != "stdout msg" {
		t.Fatalf("stdout 期望 ['stdout msg']，得到 %v", stdoutLines)
	}

	// 验证 stderr 内容
	if len(stderrLines) != 1 || stderrLines[0] != "stderr msg" {
		t.Fatalf("stderr 期望 ['stderr msg']，得到 %v", stderrLines)
	}
}
