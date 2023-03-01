// Package cmd /*
package cmd

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: `get something in the block_chain chain`,
	Long:  `get something in the block_chain chain`,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.AddCommand(balanceCmd)
}
