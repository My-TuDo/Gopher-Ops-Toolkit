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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// 获取参数数值
		target := viper.GetString("health.target")
		timeout := viper.GetDuration("health.timeout")
		probeType := viper.GetString("health.type")

		// 防御性编程：如果命令行没传，同时配置文件也没有写target，才报错提示用户
		if target == "" {
			fmt.Println("错误：未指定探测目标。请通过 --target 参数或在配置文件中设置！")
			os.Exit(1)
		}

		// 分支一：单项协议探测
		if p, exists := prober.Probers[probeType]; exists {
			res := p.Probe(target, timeout)

			fmt.Println(res.String())

			if res.Status != "健康" {
				os.Exit(1)
			}
			return
		}

		// 分支二：全项探测
		if probeType == "all" {
			fmt.Println("开始全项探测...")
			allSuccess := true
			for name, p := range prober.Probers {
				// 当 name == "tcp" 时，去 config 拿 health.target.tcp
				// 当 name == "http" 时，去 config 拿 health.target.http
				// 当 name == "dns" 时，去 config 拿 health.target.dns
				targetKey := fmt.Sprintf("health.target.%s", name)
				specificTarget := viper.GetString(targetKey)

				// 如果用户在 config.yaml 中漏配置了某项目标，则给出提示并跳过或报错
				if specificTarget == "" {
					fmt.Printf("警告：检测到 [%s] 探测项，但配置文件中未找到 %s 的配置，已自动跳过该项探测！\n", name, targetKey)
					continue
				}

				// 进行探测
				res := p.Probe(specificTarget, timeout)
				fmt.Println(res.String())

				if res.Status != "健康" {
					allSuccess = false
				}
			}

			if allSuccess {
				fmt.Println("所有探测项均通过！")
			} else {
				os.Exit(1)
			}
			return
		}

		// 无效的探测类型
		fmt.Printf("错误：无效的探测类型 '%s'。支持 'tcp', 'http', 'all'。\n", probeType)
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.Flags().StringP("target", "", "localhost:3306", "服务目标地址")
	healthCmd.Flags().DurationP("timeout", "", 5*time.Second, "连接超时时间")
	healthCmd.Flags().StringP("type", "", "tcp", "探测类型（tcp|http|all）")

	// 绑定参数到 viper
	viper.BindPFlag("health.target", healthCmd.Flags().Lookup("target"))
	viper.BindPFlag("health.timeout", healthCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("health.type", healthCmd.Flags().Lookup("type"))
}
