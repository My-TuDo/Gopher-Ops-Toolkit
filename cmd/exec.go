/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/executor"
	"github.com/spf13/cobra"
)

// execCmd represents the exec command
var execCmd = &cobra.Command{
	Use:   "exec [command...]",
	Short: "在本地执行 shell 命令",
	Long: `在本地执行 shell 命令，支持流式输出实时查看结果。

示例：
  gopt exec "go version"
  gopt exec --timeout 5s "sleep 10"
  gopt exec --stream "ping -c 4 baidu.com"
  gopt exec --output json "df -h"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("需要指定要执行的命令")
		}
		command := args[0]

		timeout, _ := cmd.Flags().GetDuration("timeout")
		streamMode, _ := cmd.Flags().GetBool("stream")
		outputFormat, _ := cmd.Flags().GetString("output")

		// 创建带超时的 Context
		ctx := context.Background()
		if timeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		// 监听系统信号（Ctrl+C），优雅退出
		sigCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
		defer stop()

		// 创建执行器并运行
		exec := executor.NewLocalExecutor()
		resultCh, err := exec.Exec(sigCtx, command)
		if err != nil {
			return fmt.Errorf("启动命令失败: %w", err)
		}

		// 根据输出格式选择输出方式
		switch outputFormat {
		case "json":
			return outputJSON(resultCh)
		default:
			return outputStream(streamMode, resultCh)
		}
	},
}

func init() {
	rootCmd.AddCommand(execCmd)

	execCmd.Flags().DurationP("timeout", "t", 0, "命令超时时间（如 5s、1m）")
	execCmd.Flags().BoolP("stream", "s", false, "流式输出模式（显示行号和时间戳）")
	execCmd.Flags().StringP("output", "o", "default", "输出格式：default 或 json")
}

// outputStream 默认输出格式 — 逐行打印
func outputStream(streamMode bool, resultCh <-chan executor.ExecResult) error {
	lineNum := 0
	for res := range resultCh {
		if res.Done {
			if res.Err != nil {
				fmt.Fprintf(os.Stderr, "\n命令退出（code=%d）: %v\n", res.ExitCode, res.Err)
			} else {
				fmt.Fprintf(os.Stderr, "\n命令执行成功（code=%d）\n", res.ExitCode)
			}
			return nil
		}

		lineNum++
		prefix := ""
		if res.IsStderr {
			prefix = "[stderr] "
		}
		if streamMode {
			fmt.Printf("[%4d] %s%s\n", lineNum, prefix, res.Line)
		} else {
			if res.IsStderr {
				fmt.Fprintln(os.Stderr, res.Line)
			} else {
				fmt.Println(res.Line)
			}
		}
	}
	return nil
}

// outputJSON JSON 格式输出
func outputJSON(resultCh <-chan executor.ExecResult) error {
	// 先收集所有行，再一次性输出 JSON（因为 JSON 是结构化格式，不适合流式）
	var lines []executor.ExecResult
	for res := range resultCh {
		lines = append(lines, res)
	}

	if len(lines) > 0 {
		last := lines[len(lines)-1]
		if last.Done {
			if last.Err != nil {
				return last.Err
			}
		}
	}

	// 简易 JSON 输出 — 可以升级成 encoding/json
	fmt.Print("[")
	for i, res := range lines {
		if i > 0 {
			fmt.Print(",")
		}
		fmt.Printf(`{"line":%q,"stderr":%v}`, res.Line, res.IsStderr)
	}
	fmt.Println("]")
	return nil
}
