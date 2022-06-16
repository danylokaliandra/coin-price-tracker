package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

// Set configuration default values
func setDefaultValues() {
	viper.SetDefault("tracklet.maxHistory", 365)
	viper.SetDefault("tracklet.timeout", 30)
	viper.SetDefault("tracklet.retryDelay", 10)
	viper.SetDefault("tracklet.maxRetries", 12)

	viper.SetDefault("aggregators.coingecko.apiBaseUrl", "https://api.coingecko.com")

	viper.SetDefault("exchanges.binance.apiBaseUrl", "https://api.binance.com")

	viper.SetDefault("exchanges.kucoin.apiBaseUrl", "https://api.kucoin.com")
}

// Load configuration file
func LoadConfig() error {
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.tracklet")
	viper.SetConfigName("tracklet")
	viper.SetConfigType("yaml")

	setDefaultValues()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("could not retrieve config file: %v", err)
	}

	return nil
}
