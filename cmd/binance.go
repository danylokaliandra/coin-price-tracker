package cmd

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/tracklet/pkg/binance"
	"github.com/eliasbokreta/tracklet/pkg/utils"
	"github.com/spf13/cobra"
)

var cmdBinance = &cobra.Command{
	Use:   "binance",
	Short: "Deal with binance",
}

var cmdBinanceProcess = &cobra.Command{
	Use:   "process",
	Short: "Process binance data",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.NewConfig()
		if err := config.LoadConfig("./config"); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		binance := binance.NewBinance(config.Exchanges["binance"], "https://api.binance.com")
		binance.ProcessBinanceData()
	},
}

func binanceCmdInit() {
	rootCmd.AddCommand(cmdBinance)
	cmdBinance.AddCommand(cmdBinanceProcess)
}
