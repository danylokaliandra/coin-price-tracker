package coingecko

import (
	"encoding/json"
	"fmt"
)

const (
	coinListEndpoint   = "/api/v3/coins/list"
	coinPricesEndpoint = "/api/v3/coins"
)

type CoinList struct {
	Coins []struct {
		ID     string `json:"id"`
		Symbol string `json:"symbol"`
		Name   string `json:"name"`
	}
}

func GetCoinList() (*CoinList, error) {
	client := NewClient()

	body, err := client.RequestWithRetries(coinListEndpoint, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("could not request coin list endpoint: %w", err)
	}

	coinList := CoinList{}
	if err := json.Unmarshal(body, &coinList.Coins); err != nil {
		return nil, fmt.Errorf("could not unmarshal coin list: %w", err)
	}

	return &coinList, nil
}

type CoinPrice struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	MarketData struct {
		CurrentPrice struct {
			EUR float64 `json:"eur"`
			USD float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

func GetCoinPrice(id string) (*CoinPrice, error) {
	client := NewClient()
	params := map[string]string{
		"tickers":        "false",
		"market_data":    "true",
		"community_data": "false",
		"developer_data": "false",
	}
	body, err := client.RequestWithRetries(fmt.Sprintf("%s/%s", coinPricesEndpoint, id), params)
	if err != nil {
		return nil, fmt.Errorf("could not request coin price endpoint: %w", err)
	}

	coinPrice := CoinPrice{}
	if err := json.Unmarshal(body, &coinPrice); err != nil {
		return nil, fmt.Errorf("could not unmarshal coin price: %w", err)
	}

	return &coinPrice, nil
}
