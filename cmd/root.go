package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var rootCmd = &cobra.Command{Use: "tracklet"}

func initCmd() {
	cobra.OnInitialize()
	binanceCmdInit()
	kucoinCmdInit()
}

func Execute() error {
	initCmd()

	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("got error %v", err)
	}

	return nil
}
