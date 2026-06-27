/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Gopher Ops Toolkit is a collection of tools and utilities designed to assist Go developers in various aspects of their development workflow. It provides a set of commands and features that can help with tasks such as code generation, project scaffolding, dependency management, and more. The toolkit aims to streamline the development process and enhance productivity for Go developers.`,
	Run: func(cmd *cobra.Command, args []string) {
		verbose := viper.GetBool("verbose")
		version := "v0.2.0"
		if verbose {
			fmt.Printf("version: %s\nconfig: %s\n", version, cfgFile)
		} else {
			fmt.Println(version)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
