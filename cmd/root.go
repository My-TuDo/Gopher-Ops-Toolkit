/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "Gopher-Ops-Toolkit",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/config.yaml", "config file (default is $HOME/.config.yaml)")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// 输出参数
	rootCmd.PersistentFlags().StringP("output", "o", "table", "输出格式：table 或 json")
	rootCmd.PersistentFlags().BoolP("save", "", false, "是否保存输出结果到文件")
	rootCmd.PersistentFlags().StringP("log-dir", "", "./logs", "日志文件保存目录")

	// 可从配置文件中读取参数
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("save", rootCmd.PersistentFlags().Lookup("save"))
	viper.BindPFlag("log-dir", rootCmd.PersistentFlags().Lookup("log-dir"))
}

func initConfig() {
	if cfgFile != "" {
		// 用户显式传了 --config 参数，就使用路径
		viper.SetConfigFile(cfgFile)
	} else {
		// 没传，就使用系统默认路径
		viper.SetConfigName("config")  // 配置文件名（不带扩展名）
		viper.SetConfigType("yaml")    // 文件类型
		viper.AddConfigPath("configs") // 添加配置文件所在的路径
		viper.AddConfigPath(".")       // 也可以在当前目录下寻找配置文件
		viper.AddConfigPath("$HOME")   // 也可以在用户的 home 目录下寻找配置文件
	}

	viper.AutomaticEnv() // 读取环境变量 环境变量也纳入优先级

	if err := viper.ReadInConfig(); err != nil {
		// 读不到配置文件，可能是用户想用默认值
		// 只有传了 --config 但是读不到的时候才报错
		if cfgFile != "" {
			panic(fmt.Errorf("config file not found:%s", cfgFile))
		}
	}
}
