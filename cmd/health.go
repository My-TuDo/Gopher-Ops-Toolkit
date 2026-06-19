/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
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

		// ———————————————— --all 模式 ————————————————
		if allMode {
			fmt.Println("执行全部探测...")

			allSuccess := true
			for name, p := range prober.Probers {
				targetKey := fmt.Sprintf("health.%s.target", name)
				specificTarget := viper.GetString(targetKey)

				if specificTarget == "" {
					fmt.Printf("跳过 %s 探测，未配置目标\n", name)
					continue
				}
				res := p.Probe(specificTarget, timeout)
				fmt.Println(res.String())
				if res.Status != "健康" {
					allSuccess = false
				}
			}
			if !allSuccess {
				os.Exit(1)
			}
			fmt.Println("所有探测完成，服务状态正常")
			return
		}

		// ———————————————— 单项探测 ————————————————
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)

	// 通用参数
	healthCmd.Flags().BoolP("all", "a", false, "执行全部探测（从配置文件读取）")
	healthCmd.Flags().StringP("type", "t", "tcp", "指定探测类型（tcp|http|dns)")
	healthCmd.Flags().DurationP("timeout", "", 5*time.Second, "探测超时时间`")

	// TCP 参数
	healthCmd.Flags().StringP("host", "", "", "TCP 探测的主机地址")
	healthCmd.Flags().StringP("port", "", "", "TCP 探测的端口号")

	// HTTP 参数
	healthCmd.Flags().StringP("url", "", "", "HTTP 探测的目标 URL")
	healthCmd.Flags().StringP("method", "X", "GET", "HTTP 探测使用的请求方法")
	healthCmd.Flags().StringArrayP("header", "H", nil, "HTTP 请求头")

	// DNS 参数
	healthCmd.Flags().StringP("domain", "", "", "DNS 探测的域名")
	healthCmd.Flags().StringP("record-type", "", "A", "DNS 记录类型（A/AAAA/CNAME/MX/TXT）")

	// 绑定参数到 viper
	viper.BindPFlag("health.timeout", healthCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("health.type", healthCmd.Flags().Lookup("type"))
}
