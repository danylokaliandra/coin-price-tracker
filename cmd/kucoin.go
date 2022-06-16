package cmd

import (
	"github.com/eliasbokreta/tracklet/pkg/kucoin"
	"github.com/spf13/cobra"
)

var cmdKucoin = &cobra.Command{
	Use:   "kucoin",
	Short: "Deal with Kucoin",
}

var cmdKucoinProcess = &cobra.Command{
	Use:   "process",
	Short: "Process Kucoin data",
	Run: func(cmd *cobra.Command, args []string) {
		kucoin := kucoin.New()
		kucoin.ProcessKucoinData(verbose)
	},
}

func kucoinCmdInit() {
	rootCmd.AddCommand(cmdKucoin)

	cmdKucoin.AddCommand(cmdKucoinProcess)
	cmdKucoinProcess.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print json data while processing")
}
