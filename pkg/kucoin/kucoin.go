package kucoin

import (
	"fmt"
	"os"

	"github.com/eliasbokreta/tracklet/pkg/utils"
	log "github.com/sirupsen/logrus"
)

type Kucoin struct {
	Accounts        *Accounts
	DepositHistory  *DepositHistory
	WithdrawHistory *WithdrawHistory
}

// Create a new Kucoin object
func New() *Kucoin {
	return &Kucoin{}
}

// Retrieve all account data from Kucoin
func (k *Kucoin) ProcessKucoinData(verbose bool) {
	log.Info("Starting process Kucoin data...")

	// EXCHANGE'S ACCOUNTS
	log.Info("Fetching accounts data...")
	accounts, err := GetAccounts()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k.Accounts = accounts

	if verbose {
		if err := utils.OutputResult(k.Accounts); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// DEPOSIT HISTORY
	log.Info("Fetching deposit history data...")
	depositHistory, err := GetDepositHistory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k.DepositHistory = depositHistory

	if verbose {
		if err := utils.OutputResult(k.DepositHistory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// WITHDRAW HISTORY
	log.Info("Fetching withdraw history data...")
	withdrawHistory, err := GetWithdrawHistory()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	k.WithdrawHistory = withdrawHistory

	if verbose {
		if err := utils.OutputResult(k.WithdrawHistory); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}
