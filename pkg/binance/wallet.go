// Handles data calculation logic
package binance

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/eliasbokreta/tracklet/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Wallet struct {
	Holdings map[string]Holdings `json:"holdings"`
	Stats    Stats               `json:"stats"`
}

type Holdings struct {
	Quantity     float64 `json:"quantity"`
	CurrentValue float64 `json:"currentValue"`
}

type Stats struct {
	TotalInvested float64 `json:"totalInvested"`
	TotalValue    float64 `json:"totalValue"`
	GainValue     float64 `json:"gainValue"`
	TotalAssets   int     `json:"totalAssets"`
}

// Create a new Wallet object
func NewWallet() *Wallet {
	return &Wallet{
		Holdings: make(map[string]Holdings),
		Stats: Stats{
			TotalInvested: 0,
			TotalValue:    0,
			GainValue:     0,
			TotalAssets:   0,
		},
	}
}

// Calculate all fiat payments
func (w *Wallet) calculateFiatPayments() {
	log.Info("Calculating fiat payments...")

	fiatPayments := FiatPayments{}
	data := utils.LoadFromFile("fiat_payments.json")

	if err := json.Unmarshal(data, &fiatPayments); err != nil {
		log.Errorf("Could not unmarshall data: %v", err)
		return
	}

	for _, fp := range fiatPayments.Data {
		sourceAmount, err := strconv.ParseFloat(fp.SourceAmount, 64)
		if err != nil {
			log.Errorf("Could not convert string to float: %v", err)
			return
		}

		if fp.Status == "Completed" {
			w.Stats.TotalInvested += sourceAmount

			obtainAmount, err := strconv.ParseFloat(fp.ObtainAmount, 64)
			if err != nil {
				log.Errorf("Could not convert string to float: %v", err)
				return
			}

			previousTotal := w.Holdings[fp.CryptoCurrency].Quantity
			w.Holdings[fp.CryptoCurrency] = Holdings{
				Quantity: previousTotal + obtainAmount,
			}
		}
	}
}

// Calculate all trades
func (w *Wallet) calculateTrades() {
	log.Info("Calculating trades...")

	tradingHistory := []TradingHistory{}
	data := utils.LoadFromFile("trading_history.json")

	if err := json.Unmarshal(data, &tradingHistory); err != nil {
		log.Errorf("Could not unmarshall data: %v", err)
		return
	}

	tradingPairs := TradingPairs{}
	data = utils.LoadFromFile("trading_pairs.json")

	if err := json.Unmarshal(data, &tradingPairs); err != nil {
		log.Errorf("Could not unmarshall data: %v", err)
		return
	}

	for _, th := range tradingHistory {
		// Allow to retrieve separate assets from a whole symbol
		for _, tp := range tradingPairs.Symbols {
			if th.Symbol == tp.Symbol {
				th.BaseAsset = tp.BaseAsset
				th.QuoteAsset = tp.QuoteAsset
			}
		}

		if th.IsBuyer {
			// Base asset bought
			obtainAmount, err := strconv.ParseFloat(th.Quantity, 64)
			if err != nil {
				log.Errorf("Could not convert string to float: %v", err)
				return
			}

			previousTotal := w.Holdings[th.BaseAsset].Quantity
			w.Holdings[th.BaseAsset] = Holdings{
				Quantity: previousTotal + obtainAmount,
			}

			// Currency used to buy
			obtainAmount, err = strconv.ParseFloat(th.QuoteQuantity, 64)
			if err != nil {
				log.Errorf("Could not convert string to float: %v", err)
				return
			}

			previousTotal = w.Holdings[th.QuoteAsset].Quantity
			w.Holdings[th.QuoteAsset] = Holdings{
				Quantity: previousTotal - obtainAmount,
			}
		} else {
			// Base asset sold
			obtainAmount, err := strconv.ParseFloat(th.Quantity, 64)
			if err != nil {
				log.Errorf("Could not convert string to float: %v", err)
				return
			}

			previousTotal := w.Holdings[th.BaseAsset].Quantity
			w.Holdings[th.BaseAsset] = Holdings{
				Quantity: previousTotal - obtainAmount,
			}

			// Currency got
			obtainAmount, err = strconv.ParseFloat(th.QuoteQuantity, 64)
			if err != nil {
				log.Errorf("Could not convert string to float: %v", err)
				return
			}

			previousTotal = w.Holdings[th.QuoteAsset].Quantity
			w.Holdings[th.QuoteAsset] = Holdings{
				Quantity: previousTotal + obtainAmount,
			}
		}
	}
}

// Retrieve prices for all assets
func (w *Wallet) calculatePrices() {
	log.Info("Calculating prices...")
	client := NewClient()
	for asset, d := range w.Holdings {
		symbol := fmt.Sprintf("%sUSDT", asset)
		log.Info(symbol)
		tickerPrice, err := GetTickerPrice(client, symbol)
		if err != nil {
			log.Errorf("Could not get ticker price: %v", err)
		}

		price, err := strconv.ParseFloat(tickerPrice.Price, 64)
		if err != nil {
			log.Errorf("Could not convert string to float: %v", err)
			return
		}

		total := price * w.Holdings[asset].Quantity
		d.CurrentValue = total
		w.Holdings[asset] = d
	}
}

// Calculate global wallet stats
func (w *Wallet) calculateStats() {
	log.Info("Calculating wallet stats...")
	for _, asset := range w.Holdings {
		w.Stats.TotalAssets += 1
		w.Stats.TotalValue += asset.CurrentValue
	}

	w.Stats.GainValue = (w.Stats.TotalValue - w.Stats.TotalInvested)
}

// Process Binance fetched data
func (w *Wallet) ProcessWallet() {
	w.calculateFiatPayments()
	w.calculateTrades()
	w.calculatePrices()
	w.calculateStats()

	if err := utils.OutputResult(w); err != nil {
		log.Errorf("Could not output result: %v", err)
		return
	}
}
