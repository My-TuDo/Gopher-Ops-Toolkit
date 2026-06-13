/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

var (
	host    string
	port    int
	timeout time.Duration
)

type Result struct {
	Name    string `json:"name"`              // 项目名称
	Target  string `json:"target"`            // 目标地址
	Status  string `json:"status"`            // 健康状态
	Error   string `json:"error,omitempty"`   // 错误信息（如果有）
	Latency int64  `json:"latency,omitempty"` // 响应时间（如果健康）
}

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "检查服务的健康状态",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()
		result := checkTCP(host, port)
		if result == nil {
			latency := time.Since(start).Milliseconds()
			fmt.Printf("服务 %s:%d 健康 (响应时间: %dms)\n", host, port, latency)
		} else {
			fmt.Printf("服务 %s:%d 不健康\n", host, port)
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.Flags().StringVarP(&host, "host", "", "localhost", "服务主机地址")
	healthCmd.Flags().IntVarP(&port, "port", "", 3306, "服务端口")
	healthCmd.Flags().DurationVarP(&timeout, "timeout", "", 5*time.Second, "连接超时时间")
}

// Sting 格式化输出结果
func (r Result) String() string {
	if r.Error != "" {
		return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s, 错误: %s", r.Name, r.Target, r.Status, r.Error)
	}
	return fmt.Sprintf("项目: %s, 目标: %s, 状态: %s", r.Name, r.Target, r.Status)
}

// TCP 端口探测，
func checkTCP(host string, port int) error {
	_, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	return err
}

// HTTP 探测
func checkHTTP(url string) error {
	// 实现HTTP探测逻辑

	// 发送HTTP请求并检查响应
	client := http.Client{}

	// 设置动态超时控制
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("创建HTTP请求失败: %v", err)
	}

	// 发送请求并检查响应
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("发送HTTP请求失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	// 检查HTTP状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil // 健康
	}
	return fmt.Errorf("HTTP状态码异常: %d", resp.StatusCode)
}
