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
		// 不传入参数时，默认进行 `all` 探测

		// 获得基础全局参数
		timeout := viper.GetDuration("health.timeout")
		probeType := viper.GetString("health.type")

		// 用户传递 type==all or type=="" 时，执行所有探测
		if probeType == "all" || probeType == "" {
			fmt.Println("执行所有探测...")
			allSuccess := true
			for name, p := range prober.Probers {
				// 动态匹配
				targetKey := fmt.Sprintf("health.targets.%s", name)
				specificTarget := viper.GetString(targetKey)

				// 防御性编程：如果用户没有为某个探测类型传递目标地址，警告并跳过。
				if specificTarget == "" {
					fmt.Printf("警告：检测到 [%s] 探测项，但未在配置文件中找到 %s 的配置，已自动跳过\n", name, targetKey)
					continue
				}

				// 派发专属 target 拨测
				res := p.Probe(specificTarget, timeout)
				fmt.Println(res.String())

				if res.Status != "健康" {
					allSuccess = false
				}
			}
			if !allSuccess {
				os.Exit(1)
			}
			fmt.Println("所有探测项均通过！")
			return
		}

		// 用户传递单项探测
		p, exists := prober.Probers[probeType]
		if !exists {
			fmt.Printf("错误：无效的探测类型 '%s'。支持 'tcp', 'http', 'all'。\n", probeType)
			os.Exit(1)
		}

		// 防御性编程：如果用户没有传递target。
		target := viper.GetString("health.target")
		if target == "" {
			fmt.Printf("错误：未指定目标地址。请使用 --target 参数指定目标地址。\n")
			os.Exit(1)
		}

		// 执行单项探测
		res := p.Probe(target, timeout)
		fmt.Println(res.String())
		if res.Status != "健康" {
			os.Exit(1)
		}
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
