package cmd

import (
	"github.com/eliasbokreta/tracklet/pkg/binance"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

var cmdBinance = &cobra.Command{
	Use:   "binance",
	Short: "Deal with Binance",
}

var cmdBinanceProcess = &cobra.Command{
	Use:   "process",
	Short: "Process Binance data",
	Run: func(cmd *cobra.Command, args []string) {
		binance := binance.NewBinance()
		binance.ProcessBinanceData(verbose)
	},
}

var cmdBinanceWallet = &cobra.Command{
	Use:   "wallet",
	Short: "Get binance wallet",
	Run: func(cmd *cobra.Command, args []string) {
		wallet := binance.NewWallet()
		wallet.ProcessWallet()
	},
}

func binanceCmdInit() {
	rootCmd.AddCommand(cmdBinance)

	cmdBinance.AddCommand(cmdBinanceProcess)
	cmdBinanceProcess.Flags().BoolVarP(&verbose, "verbose", "v", false, "Print json data while processing")

	cmdBinance.AddCommand(cmdBinanceWallet)
}
