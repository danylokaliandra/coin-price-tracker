package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{Use: "tracklet"}

func initCmd() {
	cobra.OnInitialize()
	binanceCmdInit()
}

func Execute() error {
	initCmd()
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("got error %s", err.Error())
	}

	return nil
}
