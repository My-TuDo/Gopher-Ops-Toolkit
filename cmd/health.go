/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/output"
	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/prober"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// healthCmd represents the health command
var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "检查服务的健康状态",
	Run: func(cmd *cobra.Command, args []string) {
		// 变量声明
		rounds := viper.GetInt("health.rounds")
		timeout := viper.GetDuration("health.timeout")
		probeType := viper.GetString("health.type")
		allMode, _ := cmd.Flags().GetBool("all")
		var results []prober.Result               // 用于存储探测结果
		outputFormat := viper.GetString("output") // 获取输出格式
		save := viper.GetBool("save")             // 获取是否保存结果到文件
		logDir := viper.GetString("log-dir")      // 获取日志目录

		// ———————————————— --all 模式 ————————————————
		if allMode {
			if outputFormat != "json" {
				fmt.Println("执行全部探测...")
			}

			allSuccess := true
			for name, p := range prober.Probers {
				targetKey := fmt.Sprintf("health.targets.%s", name)
				specificTarget := viper.GetString(targetKey)

				if specificTarget == "" {
					fmt.Printf("跳过 %s 探测，未配置目标\n", name)
					continue
				}
				res := prober.MultiRoundProbe(p, specificTarget, timeout, rounds)
				results = append(results, res)
				if res.Status != "健康" {
					allSuccess = false
				}
			}

			// 输出表格
			output.RenderResult(results, outputFormat)

			// 保存结果到文件
			if save {
				if err := output.SaveToFile(results, logDir); err != nil {
					fmt.Fprintf(os.Stderr, "保存日志失败: %v\n", err)
				}
			}

			if !allSuccess {
				os.Exit(1)
			}
			if allSuccess && outputFormat != "json" {
				fmt.Println("所有探测完成，服务状态正常")
			}
			return
		}

		// ———————————————— 单项探测 ————————————————
		// 根据类型组装 target
		var probeInstance prober.Prober // 用于存储带参数的探针实例（如 HTTPProber）
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
			probeInstance = prober.Probers["tcp"]

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
			keyword, _ := cmd.Flags().GetString("keyword")
			if url == "" {
				fmt.Println("[ERROR] HTTP 探测需要 --url 参数")
				os.Exit(1)
			}

			// 将 []string 解析为 map[string]string
			headerMap := make(map[string]string)
			for _, h := range headers {
				if k, v, found := strings.Cut(h, ":"); found {
					headerMap[strings.TrimSpace(k)] = strings.TrimSpace(v)
				}
			}
			target = url
			// 构造带参数的 HTTPProber，不走 Probers 注册表
			probeInstance = &prober.HTTPProber{
				Method:  method,
				Headers: headerMap,
				Keyword: keyword,
			}

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
			dnsServer, _ := cmd.Flags().GetString("dns-server")
			if domain == "" {
				fmt.Println("[ERROR] DNS 探测需要 --domain 参数")
				os.Exit(1)
			}

			target = domain
			probeInstance = &prober.DNSProber{
				RecordType: recordType,
				DNSServer:  dnsServer,
			}

		default:
			fmt.Printf("[ERROR] 不支持的探测类型: '%s',支持的类型有: tcp, http, dns\n", probeType)
			os.Exit(1)
		}

		res := prober.MultiRoundProbe(probeInstance, target, timeout, rounds)
		results = append(results, res)

		// 输出格式选择，默认为 table
		output.RenderResult(results, outputFormat)

		// 保存结果到文件
		if save {
			if err := output.SaveToFile(results, logDir); err != nil {
				fmt.Fprintf(os.Stderr, "保存日志失败: %v\n", err)
			}
		}

		if res.Status != "健康" {
			os.Exit(1)
		}
		if outputFormat != "json" {
			fmt.Println("探测完成，服务状态正常")
		}
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)

	// 通用参数
	healthCmd.Flags().BoolP("all", "a", false, "执行全部探测（从配置文件读取）")
	healthCmd.Flags().StringP("type", "t", "tcp", "指定探测类型(tcp|http|dns)")
	healthCmd.Flags().DurationP("timeout", "", 5*time.Second, "探测超时时间")
	// 新增 --rounds 参数
	healthCmd.Flags().IntP("rounds", "r", 1, "探测轮数，默认为 1")

	// TCP 参数
	healthCmd.Flags().StringP("host", "", "localhost", "TCP 探测的主机地址")
	healthCmd.Flags().StringP("port", "", "8080", "TCP 探测的端口号")

	// HTTP 参数
	healthCmd.Flags().StringP("url", "", "", "HTTP 探测的目标 URL")
	healthCmd.Flags().StringP("method", "m", "GET", "HTTP 探测使用的请求方法")
	healthCmd.Flags().StringArrayP("header", "H", nil, "HTTP 请求头")
	healthCmd.Flags().StringP("keyword", "k", "", "HTTP 响应体中期望包含的关键字")

	// DNS 参数
	healthCmd.Flags().StringP("domain", "", "", "DNS 探测的域名")
	healthCmd.Flags().StringP("record-type", "", "A", "DNS 记录类型(A/AAAA/CNAME/MX/TXT)")
	healthCmd.Flags().StringP("dns-server", "", "", "DNS 服务器地址")

	// 绑定参数到 viper
	viper.BindPFlag("health.timeout", healthCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("health.type", healthCmd.Flags().Lookup("type"))
	viper.BindPFlag("health.rounds", healthCmd.Flags().Lookup("rounds"))
}
