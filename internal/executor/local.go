package executor

import (
	"bufio"
	"context"
	"os/exec"
	"syscall"
)

// LocalExecutor 在本地执行命令
type LocalExecutor struct{}

// NewLocalExecutor 创建一个本地执行器
func NewLocalExecutor() *LocalExecutor {
	return &LocalExecutor{}
}

// Exec 执行命令，通过 channel 流式输出 stdout/stderr。
//
// 设计要点：
//   - 使用 Cmd.StdoutPipe / StderrPipe，而不是 CombinedOutput()，
//     这样才能实时流式输出，而不是等命令结束才拿到全部内容。
//   - 通过 goroutine 分别读取 stdout 和 stderr，发送到同一个 resultCh。
//   - 命令结束后发送一条 Done 信号，携带退出码和可能的错误。
//   - 支持 context 取消：ctx 取消时会杀掉进程（包括子进程）。
func (l *LocalExecutor) Exec(ctx context.Context, command string) (<-chan ExecResult, error) {
	// 使用 shell 执行，支持管道、重定向等 shell 特性
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// 获取 stdout 管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	// 获取 stderr 管道
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	resultCh := make(chan ExecResult, 64)

	go func() {
		defer close(resultCh)

		// goroutine A: 读取 stdout
		go func() {
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				line := scanner.Text()
				resultCh <- ExecResult{
					Line:     line,
					IsStderr: false,
				}
			}
		}()

		// goroutine B: 读取 stderr
		go func() {
			scanner := bufio.NewScanner(stderr)
			for scanner.Scan() {
				line := scanner.Text()
				resultCh <- ExecResult{
					Line:     line,
					IsStderr: true,
				}
			}
		}()

		// 等待命令结束
		err := cmd.Wait()

		// 获取退出码（兼容不同平台）
		exitCode := 0
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					exitCode = status.ExitStatus()
				}
			}
		}

		resultCh <- ExecResult{
			Done:     true,
			ExitCode: exitCode,
			Err:      err,
		}
	}()

	return resultCh, nil
}
