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

func GetTradingPairs(client *Client) error {
	body, err := client.RequestWithRetries("/api/v1/exchangeInfo")
	if err != nil {
		return fmt.Errorf("could not request exchangeInfo endpoint")
	}

	tradingPairs := TradingPairs{}
	if err := json.Unmarshal(body, &tradingPairs); err != nil {
		return fmt.Errorf("could not unmarshal trading pairs")
	}

	output, err := json.MarshalIndent(tradingPairs, "", " ")
	if err != nil {
		return fmt.Errorf("could not marshal trading pairs")
	}
	fmt.Println(string(output))
	return nil
}
