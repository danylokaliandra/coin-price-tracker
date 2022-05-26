// Handles Binance package logic to fetch all data
package binance

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/tracklet/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Binance struct {
	TradingPairs    *TradingPairs
	FiatPayments    *FiatPayments
	TradingHistory  *[]TradingHistory
	DustConversion  *DustConversion
	DividendHistory *DividendHistory
	DepositHistory  *[]DepositHistory
	WithdrawHistory *[]WithdrawHistory
}

// Create a new Binance object
func NewBinance() *Binance {
	return &Binance{}
}

// Write all fetched data to files
func (b *Binance) saveDataToFile() {
	if err := utils.WriteToFile("trading_pairs", b.TradingPairs); err != nil {
		log.Errorf("Could not write data to trading_pairs: %v", err)
	}

	if err := utils.WriteToFile("fiat_payments", b.FiatPayments); err != nil {
		log.Errorf("Could not write data to fiat_payments: %v", err)
	}

	if err := utils.WriteToFile("trading_history", b.TradingHistory); err != nil {
		log.Errorf("Could not write data to trading_history: %v", err)
	}

	if err := utils.WriteToFile("dust_conversion", b.DustConversion); err != nil {
		log.Errorf("Could not write data to dust_conversion: %v", err)
	}

	if err := utils.WriteToFile("dividend_history", b.DividendHistory); err != nil {
		log.Errorf("Could not write data to dividend_history: %v", err)
	}

	if err := utils.WriteToFile("deposit_history", b.DepositHistory); err != nil {
		log.Errorf("Could not write data to deposit_history: %v", err)
	}

	if err := utils.WriteToFile("withdraw_history", b.WithdrawHistory); err != nil {
		log.Errorf("Could not write data to withdraw_history: %v", err)
	}
}

// Retrieve all account data from Binance
func (b *Binance) ProcessBinanceData(verbose bool) {
	log.Info("Starting process Binance data...")

	log.Info("Initializing Binance client...")
	client := NewClient()

	// EXCHANGE'S TRADING PAIRS
	log.Info("Fetching trading pairs data...")
	tradingPairs, err := GetTradingPairs(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.TradingPairs = tradingPairs

	// FIAT PAYMENTS HISTORY
	log.Info("Fetching fiat payments history data...")
	fiatPayments, err := GetFiatPaymentsHistory(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.FiatPayments = fiatPayments

	if verbose {
		if err := utils.OutputResult(b.FiatPayments); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// TRADING HISTORY
	log.Info("Fetching trading history data...")
	tradingHistory, err := GetTradingHistory(client, b.TradingPairs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.TradingHistory = tradingHistory

	if verbose {
		if err := utils.OutputResult(b.TradingHistory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// DUST CONVERSION HISTORY
	log.Info("Fetching dust conversion history data...")
	dustConversion, err := GetDustConversionHistory(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.DustConversion = dustConversion

	if verbose {
		if err := utils.OutputResult(b.DustConversion); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// DIVIDEND HISTORY
	log.Info("Fetching dividend history data...")
	dividendHistory, err := GetDividendHistory(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.DividendHistory = dividendHistory

	if verbose {
		if err := utils.OutputResult(b.DividendHistory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// DEPOSIT HISTORY
	log.Info("Fetching deposit history data...")
	depositHistory, err := GetDepositHistory(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.DepositHistory = depositHistory

	if verbose {
		if err := utils.OutputResult(b.DepositHistory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// WITHDRAW HISTORY
	log.Info("Fetching withdraw history data...")
	withdrawHistory, err := GetWithdrawHistory(client)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	b.WithdrawHistory = withdrawHistory

	if verbose {
		if err := utils.OutputResult(b.WithdrawHistory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	b.saveDataToFile()
}
