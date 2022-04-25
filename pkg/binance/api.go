package binance

import (
	"encoding/json"
	"fmt"
)

type TradingPairs struct {
	Symbols []TradingPair `json:"symbols"`
}

type TradingPair struct {
	Symbol     string `json:"symbol"`
	BaseAsset  string `json:"baseAsset"`
	QuoteAsset string `json:"quoteAsset"`
}

func GetTradingPairs(client *Client) (*TradingPairs, error) {
	body, err := client.RequestWithRetries("/api/v1/exchangeInfo", map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("could not request exchangeInfo endpoint")
	}

	tradingPairs := TradingPairs{}
	if err := json.Unmarshal(body, &tradingPairs); err != nil {
		return nil, fmt.Errorf("could not unmarshal trading pairs")
	}

	return &tradingPairs, nil
}

type TradingHistory struct {
	Symbol   string `json:"symbol"`
	Id       int64  `json:"id"`
	Price    string `json:"price"`
	Quantity string `json:"qty"`
	IsBuyer  bool   `json:"isBuyer"`
}

func GetTradingHistory(client *Client) (*[]TradingHistory, error) {
	params := map[string]string{
		"symbol": "DOTBUSD",
		"limit":  "1000",
	}
	body, err := client.RequestWithRetries("/api/v3/myTrades", params)
	if err != nil {
		return nil, fmt.Errorf("could not request myTrades endpoint")
	}

	tradingHistory := []TradingHistory{}
	if err := json.Unmarshal(body, &tradingHistory); err != nil {
		return nil, fmt.Errorf("could not unmarshal trading history")
	}

	return &tradingHistory, nil
}
