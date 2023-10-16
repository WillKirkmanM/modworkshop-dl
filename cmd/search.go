/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/WillKirkmanM/modworkshop-dl/pkg/search"

	"github.com/spf13/cobra"
)


// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "The mod to search",
	Long: `Not Implemented`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please Provide a Mod to Search For!")
			return
		}

		search.ListOfMods(args)

	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
