package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Exchanges map[string]Exchange `mapstructure:"exchanges"`
}

type Exchange struct {
	APIKey    string `mapstructure:"apiKey"`
	SecretKey string `mapstructure:"secretKey"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) LoadConfig(path string) error {
	viper.SetDefault("exchanges.binance", Exchange{})
	viper.SetDefault("exchanges.coinbase", Exchange{})
	viper.SetDefault("exchanges.kucoin", Exchange{})

	viper.AddConfigPath(path)
	viper.SetConfigName("tracklet")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("could not retrieve config file")
	}

	err = viper.Unmarshal(&c)
	if err != nil {
		return fmt.Errorf("error unmarshalling config file")
	}

	return nil
}
