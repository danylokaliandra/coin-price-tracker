// Handles Binance API endpoints logic
package binance

import (
	"encoding/json"
	"fmt"

	"github.com/eliasbokreta/tracklet/pkg/utils"
)

const (
	tradingPairsEndpoint          = "/api/v1/exchangeInfo"
	fiatPaymentsEndpoint          = "/sapi/v1/fiat/payments"
	tradingHistoryEndpoint        = "/api/v3/myTrades"
	dustConversionHistoryEndpoint = "/sapi/v1/asset/dribblet"
	dividendHistoryEndpoint       = "/sapi/v1/asset/assetDividend"
	depositHistoryEndpoint        = "/sapi/v1/capital/deposit/hisrec"
	withdrawHistoryEndpoint       = "/sapi/v1/capital/withdraw/history"
)

type TradingPairs struct {
	Symbols []struct {
		Symbol     string `json:"symbol"`
		BaseAsset  string `json:"baseAsset"`
		QuoteAsset string `json:"quoteAsset"`
	} `json:"symbols"`
}

// Get trading pairs available on Binance
func GetTradingPairs() (*TradingPairs, error) {
	client := NewClient()
	body, err := client.RequestWithRetries(tradingPairsEndpoint, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("could not request trading pairs endpoint: %v", err)
	}

	tradingPairs := TradingPairs{}
	if err := json.Unmarshal(body, &tradingPairs); err != nil {
		return nil, fmt.Errorf("could not unmarshal trading pairs: %v", err)
	}

	return &tradingPairs, nil
}

type FiatPayments struct {
	Data []struct {
		OrderNo        string `json:"orderNo"`
		SourceAmount   string `json:"sourceAmount"`
		FiatCurrency   string `json:"fiatCurrency"`
		ObtainAmount   string `json:"obtainAmount"`
		CryptoCurrency string `json:"cryptoCurrency"`
		TotalFee       string `json:"totalFee"`
		Price          string `json:"price"`
		Status         string `json:"status"`
		CreateTime     int    `json:"createTime"`
	} `json:"data"`
}

// Get Binance account fiat payments history
func GetFiatPaymentsHistory() (*FiatPayments, error) {
	client := NewClient()
	dateRanges := utils.GetDateRanges(client.MaxHistory, 15)
	fiatPayments := FiatPayments{}

	for _, dateRange := range dateRanges {
		params := map[string]string{
			"beginTime":       fmt.Sprintf("%d", dateRange.StartDate),
			"endTime":         fmt.Sprintf("%d", dateRange.EndDate),
			"transactionType": "0", // 0-buy, 1-sell
			"rows":            "500",
		}

		body, err := client.RequestWithRetries(fiatPaymentsEndpoint, params)
		if err != nil {
			return nil, fmt.Errorf("could not request fiat payments endpoint: %v", err)
		}

		fiatPaymentsRange := FiatPayments{}
		if err := json.Unmarshal(body, &fiatPaymentsRange); err != nil {
			return nil, fmt.Errorf("could not unmarshal fiat payments history: %v", err)
		}

		if len(fiatPaymentsRange.Data) > 0 {
			fiatPayments.Data = append(fiatPayments.Data, fiatPaymentsRange.Data...)
		}
	}

	return &fiatPayments, nil
}

type TradingHistory struct {
	Symbol          string `json:"symbol"`
	BaseAsset       string `json:"baseAsset,omitempty"`
	QuoteAsset      string `json:"quoteAsset,omitempty"`
	Id              int64  `json:"id"`
	Price           string `json:"price"`
	Quantity        string `json:"qty"`
	QuoteQuantity   string `json:"quoteQty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
	IsBuyer         bool   `json:"isBuyer"`
	Time            int    `json:"time"`
}

// Get Binance account trading history
func GetTradingHistory(tradingPairs *TradingPairs) (*[]TradingHistory, error) {
	client := NewClient()
	tradingHistory := []TradingHistory{}

	for _, tp := range tradingPairs.Symbols {
		params := map[string]string{
			"symbol": tp.Symbol,
			"limit":  "1000",
		}
		body, err := client.RequestWithRetries(tradingHistoryEndpoint, params)
		if err != nil {
			return nil, fmt.Errorf("could not request trades endpoint: %v", err)
		}

		tradingHistoryRange := []TradingHistory{}
		if err := json.Unmarshal(body, &tradingHistoryRange); err != nil {
			return nil, fmt.Errorf("could not unmarshal trading history: %v", err)
		}

		if len(tradingHistoryRange) > 0 {
			tradingHistory = append(tradingHistory, tradingHistoryRange...)
		}
	}

	return &tradingHistory, nil
}

type DustConversion struct {
	UserAssetDribblets []struct {
		OperateTime              int    `json:"operateTime"`
		TotalTransferedAmount    string `json:"totalTransferedAmount"`
		UserAssetDribbletDetails []struct {
			FromAsset        string `json:"fromAsset"`
			Amount           string `json:"amount"`
			TransferedAmount string `json:"transferedAmount"`
		} `json:"userAssetDribbletDetails"`
	} `json:"userAssetDribblets"`
}

// Get dust conversion history
func GetDustConversionHistory() (*DustConversion, error) {
	client := NewClient()
	dateRanges := utils.GetDateRanges(client.MaxHistory, 15)
	dustConversion := DustConversion{}

	for _, dateRange := range dateRanges {
		params := map[string]string{
			"startTime": fmt.Sprintf("%d", dateRange.StartDate),
			"endTime":   fmt.Sprintf("%d", dateRange.EndDate),
		}

		body, err := client.RequestWithRetries(dustConversionHistoryEndpoint, params)
		if err != nil {
			return nil, fmt.Errorf("could not request dust conversion history endpoint: %v", err)
		}

		dustConversionRange := DustConversion{}
		if err := json.Unmarshal(body, &dustConversionRange); err != nil {
			return nil, fmt.Errorf("could not unmarshal dust conversion history: %v", err)
		}

		if len(dustConversionRange.UserAssetDribblets) > 0 {
			dustConversion.UserAssetDribblets = append(dustConversion.UserAssetDribblets, dustConversionRange.UserAssetDribblets...)
		}
	}

	return &dustConversion, nil
}

type DividendHistory struct {
	Rows []struct {
		Amount  string `json:"amount"`
		Asset   string `json:"asset"`
		DivTime int    `json:"divTime"`
	} `json:"rows"`
}

// Get dividend (staking) rewards history
func GetDividendHistory() (*DividendHistory, error) {
	client := NewClient()
	dateRanges := utils.GetDateRanges(client.MaxHistory, 15)
	dividendHistory := DividendHistory{}

	for _, dateRange := range dateRanges {
		params := map[string]string{
			"startTime": fmt.Sprintf("%d", dateRange.StartDate),
			"endTime":   fmt.Sprintf("%d", dateRange.EndDate),
			"limit":     "500",
		}

		body, err := client.RequestWithRetries(dividendHistoryEndpoint, params)
		if err != nil {
			return nil, fmt.Errorf("could not request dividend history endpoint: %v", err)
		}

		dividendHistoryRange := DividendHistory{}
		if err := json.Unmarshal(body, &dividendHistoryRange); err != nil {
			return nil, fmt.Errorf("could not unmarshal dividend history: %v", err)
		}

		if len(dividendHistoryRange.Rows) > 0 {
			dividendHistory.Rows = append(dividendHistory.Rows, dividendHistoryRange.Rows...)
		}
	}

	return &dividendHistory, nil
}

type DepositHistory struct {
	Amount     string `json:"amount"`
	Coin       string `json:"coin"`
	InsertTime int    `json:"insertTime"`
}

// Get deposit history
func GetDepositHistory() (*[]DepositHistory, error) {
	client := NewClient()
	dateRanges := utils.GetDateRanges(client.MaxHistory, 15)
	depositHistory := []DepositHistory{}

	for _, dateRange := range dateRanges {
		params := map[string]string{
			"status":    "1", // 0:pending,6: credited but cannot withdraw, 1:success
			"startTime": fmt.Sprintf("%d", dateRange.StartDate),
			"endTime":   fmt.Sprintf("%d", dateRange.EndDate),
			"limit":     "1000",
		}

		body, err := client.RequestWithRetries(depositHistoryEndpoint, params)
		if err != nil {
			return nil, fmt.Errorf("could not request deposit history endpoint: %v", err)
		}

		depositHistoryRange := []DepositHistory{}
		if err := json.Unmarshal(body, &depositHistoryRange); err != nil {
			return nil, fmt.Errorf("could not unmarshal deposit history: %v", err)
		}

		if len(depositHistoryRange) > 0 {
			depositHistory = append(depositHistory, depositHistoryRange...)
		}
	}

	return &depositHistory, nil
}

type WithdrawHistory struct {
	Amount    string `json:"amount"`
	Coin      string `json:"coin"`
	ApplyTime string `json:"applyTime"`
}

// Get withdraw history
func GetWithdrawHistory() (*[]WithdrawHistory, error) {
	client := NewClient()
	dateRanges := utils.GetDateRanges(client.MaxHistory, 15)
	withdrawHistory := []WithdrawHistory{}

	for _, dateRange := range dateRanges {
		params := map[string]string{
			"status":    "6", // 0:Email Sent,1:Cancelled 2:Awaiting Approval 3:Rejected 4:Processing 5:Failure 6:Completed
			"startTime": fmt.Sprintf("%d", dateRange.StartDate),
			"endTime":   fmt.Sprintf("%d", dateRange.EndDate),
			"limit":     "1000",
		}

		body, err := client.RequestWithRetries(withdrawHistoryEndpoint, params)
		if err != nil {
			return nil, fmt.Errorf("could not request withdraw history endpoint: %v", err)
		}

		withdrawHistoryRange := []WithdrawHistory{}
		if err := json.Unmarshal(body, &withdrawHistoryRange); err != nil {
			return nil, fmt.Errorf("could not unmarshal withdraw history: %v", err)
		}

		if len(withdrawHistoryRange) > 0 {
			withdrawHistory = append(withdrawHistory, withdrawHistoryRange...)
		}
	}

	return &withdrawHistory, nil
}
