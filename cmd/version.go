/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "View the Current Version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v2.0.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
