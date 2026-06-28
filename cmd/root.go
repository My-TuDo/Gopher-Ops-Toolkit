/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/My-TuDo/Gopher-Ops-Toolkit/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:           "gopt",
		Short:         "Gopher-Ops-Toolkit: 一个运维工具箱",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig) // 告诉 cobra 在执行任何命令之前先运行 initConfig 函数
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ~/.config/gopt/config.yaml)")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 输出参数
	rootCmd.PersistentFlags().StringP("output", "o", "table", "输出格式：table 或 json")
	rootCmd.PersistentFlags().BoolP("save", "", false, "是否保存输出结果到文件")
	rootCmd.PersistentFlags().StringP("log-dir", "", config.DefaultLogDir(), "日志文件保存目录")

	// 可从配置文件中读取参数
	_ = viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	_ = viper.BindPFlag("save", rootCmd.PersistentFlags().Lookup("save"))
	_ = viper.BindPFlag("log-dir", rootCmd.PersistentFlags().Lookup("log-dir"))
}

// initConfig 按优先级查找并加载配置文件：
//
//	1. --config 显式指定
//	2. GOPT_CONFIG 环境变量
//	3. XDG 标准路径 ~/.config/gopt/config.yaml（首次运行自动创建）
//	4. 本地开发路径 configs/config.yaml
func initConfig() {
	path, exists := config.ResolveConfigPath(cfgFile)

	if !exists {
		// 首次运行，在 XDG 路径自动创建默认配置文件
		created, err := config.EnsureDefaultConfig(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "⚠️ 警告: 无法创建默认配置文件 %s: %v\n", path, err)
			fmt.Fprintf(os.Stderr, "   内嵌默认配置足以支撑运行，如需自定义请手动创建。\n\n")
		} else if created {
			fmt.Fprintf(os.Stderr, "\n📄 默认配置文件已创建: %s\n", path)
			fmt.Fprintf(os.Stderr, "   编辑此文件可自定义探测目标、输出格式等。\n\n")
		}
	}

	// 设置加载路径
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	// 环境变量覆盖（如 GOPT_OUTPUT=json）
	viper.SetEnvPrefix("GOPT")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		// 读不到配置文件，可能是用户想用默认值
		// 只有传了 --config 但是读不到的时候才报错
		if cfgFile != "" {
			panic(fmt.Errorf("config file not found: %s", cfgFile))
		}
	}
}
