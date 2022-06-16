// Handles HTTP requests logic
package kucoin

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Client struct {
	HTTPClient *http.Client
	BaseUrl    string
	APIKey     string
	SecretKey  string
	Passphrase string
	RetryDelay time.Duration
	MaxRetries int
	MaxHistory int
}

// Create a new Client object
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * viper.GetDuration("tracklet.timeout"),
		},
		BaseUrl:    viper.GetString("exchanges.kucoin.apiBaseUrl"),
		APIKey:     viper.GetString("exchanges.kucoin.apiKey"),
		SecretKey:  viper.GetString("exchanges.kucoin.secretKey"),
		Passphrase: viper.GetString("exchanges.kucoin.passphrase"),
		RetryDelay: viper.GetDuration("tracklet.retryDelay"),
		MaxRetries: viper.GetInt("tracklet.maxRetries"),
		MaxHistory: viper.GetInt("tracklet.maxHistory"),
	}
}

// HTTP request Kucoin server time
func (c *Client) getServerTime() int64 {
	t := time.Now()

	return int64(time.Nanosecond) * t.UnixNano() / int64(time.Millisecond)
}

// Generate a HMAC signed query string for authorizing API requests
func (c *Client) generateSignature(queryString string) (string, error) {
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	_, err := mac.Write([]byte(queryString))
	if err != nil {
		return "", fmt.Errorf("could not generate hmac: %v", err)
	}
	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

// HTTP get request for a given Kucoin API endpoint
func (c *Client) request(endpoint string, parameters map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseUrl, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("error while executing request: %v", err)
	}

	serverTime := c.getServerTime()
	if err != nil {
		return nil, fmt.Errorf("error while getting server time: %v", err)
	}

	signature, err := c.generateSignature(fmt.Sprintf("%d%s%s", serverTime, req.Method, req.URL.Path))
	if err != nil {
		return nil, fmt.Errorf("error while generating signature")
	}

	passphrase, err := c.generateSignature(c.Passphrase)
	if err != nil {
		return nil, fmt.Errorf("error while generating passphrase signature")
	}

	req.Header.Add("KC-API-KEY", c.APIKey)
	req.Header.Add("KC-API-SIGN", signature)
	req.Header.Add("KC-API-TIMESTAMP", fmt.Sprint(serverTime))
	req.Header.Add("KC-API-PASSPHRASE", passphrase)
	req.Header.Add("KC-API-KEY-VERSION", "2")

	response, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while executing request: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non valid HTTP status code : %s", response.Status)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body: %v", err)
	}

	return body, nil
}

// Retry mechanism for HTTP requests
func (c *Client) RequestWithRetries(endpoint string, parameters map[string]string) ([]byte, error) {
	var body []byte
	var err error
	retries := 0

	for {
		body, err = c.request(endpoint, parameters)
		if body != nil {
			break
		}

		if retries == c.MaxRetries {
			log.Warn("Max retries exceeded...")
			break
		}

		retries += 1
		log.Warnf("Retrying... [%d/%d]: %v", retries, c.MaxRetries, err)
		time.Sleep(time.Second * c.RetryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("could not request '%s' endpoint: %v", endpoint, err)
	}

	return body, err
}
