package kucoin

import (
	"encoding/json"
	"fmt"
)

const (
	accountsEndpoint        = "/api/v1/accounts"
	depositHistoryEndpoint  = "/api/v1/deposits"
	withdrawHistoryEndpoint = "/api/v1/withdrawals"
)

type Pagination struct {
	CurrentPage int `json:"currentPage"`
	PageSize    int `json:"pageSize"`
	TotalNum    int `json:"totalNum"`
	TotalPage   int `json:"totalPage"`
}

type Accounts struct {
	Data []struct {
		ID        string `json:"id"`
		Currency  string `json:"currency"`
		Type      string `json:"type"`
		Balance   string `json:"balance"`
		Available string `json:"available"`
		Holds     string `json:"holds"`
	} `json:"data"`
}

// Get Kucoin accounts
func GetAccounts() (*Accounts, error) {
	client := NewClient()
	body, err := client.RequestWithRetries(accountsEndpoint, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("could not request accounts endpoint: %v", err)
	}

	accounts := Accounts{}
	if err := json.Unmarshal(body, &accounts); err != nil {
		return nil, fmt.Errorf("could not unmarshal accounts: %v", err)
	}

	return &accounts, nil
}

type DepositHistory struct {
	Data struct {
		Pagination Pagination
		Items      []struct {
			Address   string `json:"address"`
			Amount    string `json:"amount"`
			Fee       string `json:"fee"`
			Currency  string `json:"currency"`
			IsInner   bool   `json:"isInner"`
			Status    string `json:"status"`
			CreatedAt int64  `json:"createdAt"`
		} `json:"items"`
	} `json:"data"`
}

// Get deposit history
func GetDepositHistory() (*DepositHistory, error) {
	client := NewClient()
	body, err := client.RequestWithRetries(depositHistoryEndpoint, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("could not request deposit history endpoint: %v", err)
	}

	depositHistory := DepositHistory{}
	if err := json.Unmarshal(body, &depositHistory); err != nil {
		return nil, fmt.Errorf("could not unmarshal deposit history: %v", err)
	}

	return &depositHistory, nil
}

type WithdrawHistory struct {
	Data struct {
		Pagination Pagination
		Items      []struct {
			Address   string `json:"address"`
			Amount    string `json:"amount"`
			Fee       string `json:"fee"`
			Currency  string `json:"currency"`
			IsInner   bool   `json:"isInner"`
			Status    string `json:"status"`
			CreatedAt int64  `json:"createdAt"`
		} `json:"items"`
	} `json:"data"`
}

// Get withdraw history
func GetWithdrawHistory() (*WithdrawHistory, error) {
	client := NewClient()
	body, err := client.RequestWithRetries(withdrawHistoryEndpoint, map[string]string{})
	if err != nil {
		return nil, fmt.Errorf("could not request withdraw history endpoint: %v", err)
	}

	withdrawHistory := WithdrawHistory{}
	if err := json.Unmarshal(body, &withdrawHistory); err != nil {
		return nil, fmt.Errorf("could not unmarshal withdraw history: %v", err)
	}

	return &withdrawHistory, nil
}
