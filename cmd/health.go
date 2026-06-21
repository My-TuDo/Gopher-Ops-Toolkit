/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/prober"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "检查服务的健康状态",
	Run: func(cmd *cobra.Command, args []string) {
		timeout := viper.GetDuration("health.timeout")
		probeType, _ := cmd.Flags().GetString("type")
		allMode, _ := cmd.Flags().GetBool("all")
		rounds, _ := cmd.Flags().GetInt("rounds")

		// ———————————————— --all 模式 ————————————————
		if allMode {
			fmt.Println("执行全部探测...")

			allSuccess := true
			for name, p := range prober.Probers {
				targetKey := fmt.Sprintf("health.targets.%s", name)
				specificTarget := viper.GetString(targetKey)

				if specificTarget == "" {
					fmt.Printf("跳过 %s 探测，未配置目标\n", name)
					continue
				}
				res := prober.MultiRoundProbe(p, specificTarget, timeout, rounds)
				fmt.Println(res.String())
				if res.Status != "健康" {
					allSuccess = false
				}
			}
			if !allSuccess {
				os.Exit(1)
			}
			if allSuccess {
				fmt.Println("所有探测完成，服务状态正常")
			}
			return
		}

		// ———————————————— 单项探测 ————————————————
		p, exists := prober.Probers[probeType]
		if !exists {
			fmt.Printf("[ERROR] 无效的探测类型 '%s'。支持的类型有: tcp, http, dns\n", probeType)
			os.Exit(1)
		}

		// 根据类型组装 target
		var target string
		switch probeType {
		case "tcp":
			if cmd.Flags().Changed("url") {
				fmt.Println("[ERROR] TCP 探测不支持 --url 参数")
				os.Exit(1)
			}
			if cmd.Flags().Changed("domain") {
				fmt.Println("[ERROR] TCP 探测不支持 --domain 参数")
				os.Exit(1)
			}

			host, _ := cmd.Flags().GetString("host")
			port, _ := cmd.Flags().GetString("port")
			if host == "" || port == "" {
				fmt.Println("[ERROR] TCP 探测需要 --host 和 --port 参数")
				os.Exit(1)
			}
			target = net.JoinHostPort(host, port)

		case "http":
			if cmd.Flags().Changed("host") || cmd.Flags().Changed("port") {
				fmt.Println("[ERROR] HTTP 探测不支持 --host 或 --port 参数")
				os.Exit(1)
			}
			if cmd.Flags().Changed("domain") {
				fmt.Println("[ERROR] HTTP 探测不支持 --domain 参数")
				os.Exit(1)
			}

			url, _ := cmd.Flags().GetString("url")
			method, _ := cmd.Flags().GetString("method")
			headers, _ := cmd.Flags().GetStringArray("header")
			if url == "" {
				fmt.Println("[ERROR] HTTP 探测需要 --url 参数")
				os.Exit(1)
			}
			// 目前先拼接 url, 后续再拓展 method 和 headers 的使用
			target = url
			_ = method
			_ = headers

		case "dns":
			if cmd.Flags().Changed("host") || cmd.Flags().Changed("port") {
				fmt.Println("[ERROR] DNS 探测不支持 --host 或 --port 参数")
				os.Exit(1)
			}
			if cmd.Flags().Changed("url") {
				fmt.Println("[ERROR] DNS 探测不支持 --url 参数")
				os.Exit(1)
			}

			domain, _ := cmd.Flags().GetString("domain")
			recordType, _ := cmd.Flags().GetString("record-type")
			if domain == "" {
				fmt.Println("[ERROR] DNS 探测需要 --domain 参数")
				os.Exit(1)
			}
			// 把 recordType 拼接到 target 中，供 Probe 使用
			target = fmt.Sprintf("%s@%s", domain, recordType)
		}

		res := prober.MultiRoundProbe(p, target, timeout, rounds)
		fmt.Println(res.String())
		if res.Status != "健康" {
			os.Exit(1)
		}
		fmt.Println("探测完成，服务状态正常")
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)

	// 通用参数
	healthCmd.Flags().BoolP("all", "a", false, "执行全部探测（从配置文件读取）")
	healthCmd.Flags().StringP("type", "t", "tcp", "指定探测类型(tcp|http|dns)")
	healthCmd.Flags().DurationP("timeout", "", 5*time.Second, "探测超时时间`")
	// 新增 --rounds 参数
	healthCmd.Flags().IntP("rounds", "r", 1, "探测轮数，默认为 1")

	// TCP 参数
	healthCmd.Flags().StringP("host", "", "localhost", "TCP 探测的主机地址")
	healthCmd.Flags().StringP("port", "", "8080", "TCP 探测的端口号")

	// HTTP 参数
	healthCmd.Flags().StringP("url", "", "", "HTTP 探测的目标 URL")
	healthCmd.Flags().StringP("method", "m", "GET", "HTTP 探测使用的请求方法")
	healthCmd.Flags().StringArrayP("header", "H", nil, "HTTP 请求头")

	// DNS 参数
	healthCmd.Flags().StringP("domain", "", "", "DNS 探测的域名")
	healthCmd.Flags().StringP("record-type", "", "A", "DNS 记录类型（A/AAAA/CNAME/MX/TXT）")

	// 绑定参数到 viper
	viper.BindPFlag("health.timeout", healthCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("health.type", healthCmd.Flags().Lookup("type"))
	viper.BindPFlag("health.rounds", healthCmd.Flags().Lookup("rounds"))
}
