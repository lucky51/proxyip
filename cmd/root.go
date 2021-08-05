package cmd

import (
	"github.com/spf13/cobra"
)
var rootCmd = cobra.Command{}

func init() {
	rootCmd.AddCommand(crawlCmd)
	rootCmd.AddCommand(poolCmd)
}

func Execute()error  {
	return rootCmd.Execute()
}