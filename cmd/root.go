/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Gopher-Ops-Toolkit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	cobra.OnInitialize(initConfig) // 告诉 cobra 在执行任何命令之前先运行 initConfig 函数
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// 用户显式传了 --config 参数，就使用路径
		viper.SetConfigFile(cfgFile)
	} else {
		// 没传，就使用系统默认路径
		viper.SetConfigName(".config") // 配置文件名（不带扩展名）
		viper.SetConfigFile("yaml")    // 文件类型
		viper.AddConfigPath("$HOME/")  // 添加配置文件所在的路径
		viper.AddConfigPath(".")       // 也可以在当前目录下寻找配置文件
	}
}
