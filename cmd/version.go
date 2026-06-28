/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// 以下变量由 Makefile 通过 ldflags 注入
var (
	commitHash string
	buildTime  string
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		version := "v0.2.0"
		if verbose {
			fmt.Printf("version:   %s\n", version)
			fmt.Printf("commit:    %s\n", commitHash)
			fmt.Printf("built:     %s\n", buildTime)
			fmt.Printf("go:        %s %s/%s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
			fmt.Printf("config:    %s\n", cfgFile)
		} else {
			fmt.Println(version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolP("verbose", "v", false, "显示详细版本信息（commit、编译时间、Go 版本等）")
}
