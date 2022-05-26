// Handles data calculation logic
package binance

import (
	"encoding/json"
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
	TotalInvested  float64 `json:"totalInvested"`
	GainPercentage float64 `json:"gainPercentage"`
	GainValue      float64 `json:"gainValue"`
}

// Create a new Wallet object
func NewWallet() *Wallet {
	return &Wallet{
		Holdings: make(map[string]Holdings),
		Stats: Stats{
			TotalInvested:  0,
			GainPercentage: 0,
			GainValue:      0,
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

			// Currency gained
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

// Process Binance fetched data
func (w *Wallet) ProcessWallet() {
	w.calculateFiatPayments()
	w.calculateTrades()

	if err := utils.OutputResult(w); err != nil {
		log.Errorf("Could not output result: %v", err)
		return
	}
}
