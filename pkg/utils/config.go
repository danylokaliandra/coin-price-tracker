package utils

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Tracklet  Tracklet            `mapstructure:"tracklet"`
	Exchanges map[string]Exchange `mapstructure:"exchanges"`
}

type Tracklet struct {
	MaxHistory int           `mapstructure:"maxHistory"`
	Timeout    time.Duration `mapstructure:"timeout"`
	RetryDelay time.Duration `mapstructure:"retryDelay"`
	MaxRetries int           `mapstructure:"maxRetries"`
}

type Exchange struct {
	APIKey       string   `mapstructure:"apiKey"`
	SecretKey    string   `mapstructure:"secretKey"`
	IncludePairs []string `mapstructure:"includePairs"`
}

// Create a new Config object
func NewConfig() *Config {
	return &Config{}
}

// Set configuration default values
func setDefaultValues() {
	viper.SetDefault("tracklet.maxHistory", 365)
	viper.SetDefault("tracklet.timeout", 30)
	viper.SetDefault("tracklet.retryDelay", 5)
	viper.SetDefault("tracklet.maxRetries", 10)

	viper.SetDefault("exchanges.binance", Exchange{})
	viper.SetDefault("exchanges.binance.apiBaseUrl", "https://api.binance.com")
	viper.SetDefault("exchanges.coinbase", Exchange{})
	viper.SetDefault("exchanges.kucoin", Exchange{})
}

// Load configuration file
func (c *Config) LoadConfig() error {
	viper.AddConfigPath("./config")
	viper.AddConfigPath("$HOME/.tracklet")
	viper.SetConfigName("tracklet")
	viper.SetConfigType("yaml")

	setDefaultValues()

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("could not retrieve config file: %v", err)
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file: %v", err)
	}

	return nil
}
