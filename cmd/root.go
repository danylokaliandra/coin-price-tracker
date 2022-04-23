package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0"

var rootCmd = &cobra.Command{Use: "tracklet", Version: fmt.Sprintf("%s", version)}

func initCmd() {
	cobra.OnInitialize()
	trackletCmdInit()
}

func Execute() error {
	initCmd()
	if err := rootCmd.Execute(); err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}

	return nil
}