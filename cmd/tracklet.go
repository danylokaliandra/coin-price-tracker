package cmd

import (
	"github.com/spf13/cobra"
)

var cmdBase = &cobra.Command{
	Use:   "base",
	Short: "Deal with base",
}

func trackletCmdInit() {
	rootCmd.AddCommand(cmdBase)
}
