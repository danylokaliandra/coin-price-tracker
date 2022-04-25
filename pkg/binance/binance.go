package binance

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/tracklet/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Binance struct {
	Config  utils.Exchange
	BaseUrl string

	TradingPairs   TradingPairs
	TradingHistory []TradingHistory
}

func NewBinance(config utils.Exchange, baseUrl string) *Binance {
	return &Binance{
		Config:  config,
		BaseUrl: baseUrl,
	}
}

func (b *Binance) ProcessBinanceData() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.Println("Starting process Binance data...")

	log.Println("Initializing Binance client...")
	client := NewClient(b.BaseUrl, b.Config.APIKey, b.Config.SecretKey, 30, 5, 10)

	log.Println("Fetching trading pairs data...")
	tradingPairs, err := GetTradingPairs(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.TradingPairs = *tradingPairs

	log.Println("Fetching trading history data...")
	tradingHistory, err := GetTradingHistory(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.TradingHistory = *tradingHistory
	//utils.OutputResult(b.TradingPairs)
	if err := utils.OutputResult(b.TradingHistory); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
