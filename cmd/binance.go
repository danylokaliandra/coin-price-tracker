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

var cmdBinanceGet = &cobra.Command{
	Use:   "get",
	Short: "Get binance data",
	Run: func(cmd *cobra.Command, args []string) {
		config := utils.NewConfig()
		config.LoadConfig("./config")

		client := binance.NewClient("https://api.binance.com", config.Exchanges["binance"].APIKey, config.Exchanges["binance"].SecretKey, 30, 5, 10, 365)

		if err := binance.GetTradingPairs(client); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func binanceCmdInit() {
	rootCmd.AddCommand(cmdBinance)
	cmdBinance.AddCommand(cmdBinanceGet)
}
