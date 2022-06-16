// Handles HTTP requests logic
package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Client struct {
	HTTPClient *http.Client
	BaseUrl    string
	APIKey     string
	SecretKey  string
	RetryDelay time.Duration
	MaxRetries int
	MaxHistory int
}

type ServerTime struct {
	ServerTime int64 `json:"serverTime"`
}

// Create a new Client object
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * viper.GetDuration("tracklet.timeout"),
		},
		BaseUrl:    viper.GetString("exchanges.binance.apiBaseUrl"),
		APIKey:     viper.GetString("exchanges.binance.apiKey"),
		SecretKey:  viper.GetString("exchanges.binance.secretKey"),
		RetryDelay: viper.GetDuration("tracklet.retryDelay"),
		MaxRetries: viper.GetInt("tracklet.maxRetries"),
		MaxHistory: viper.GetInt("tracklet.maxHistory"),
	}
}

// HTTP request Binance server time
func (c *Client) getServerTime() (*ServerTime, error) {
	response, err := http.Get(fmt.Sprintf("%s/api/v3/time", c.BaseUrl))
	if err != nil {
		return nil, fmt.Errorf("could not request time endpoint: %v", err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body: %v", err)
	}

	serverTime := ServerTime{}
	if err := json.Unmarshal(body, &serverTime); err != nil {
		return nil, fmt.Errorf("could not unmarshal server time: %v", err)
	}

	return &serverTime, nil
}

// Generate a HMAC signed query string for authorizing API requests
func (c *Client) generateSignature(queryString string) (string, error) {
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	_, err := mac.Write([]byte(queryString))
	digest := mac.Sum(nil)
	if err != nil {
		return "", fmt.Errorf("could not generate hmac: %v", err)
	}
	hexDigest := hex.EncodeToString(digest)
	return hexDigest, nil
}

// Create a Binance accepted query string
func (c *Client) buildQueryString(q url.Values, params map[string]string) url.Values {
	if len(params) == 0 {
		return nil
	}
	for elem := range params {
		q.Add(elem, params[elem])
	}

	serverTime, err := c.getServerTime()
	if err != nil {
		return nil
	}
	q.Add("timestamp", fmt.Sprint(serverTime.ServerTime))

	signature, err := c.generateSignature(q.Encode())
	if err != nil {
		return nil
	}

	q.Add("signature", signature)

	return q
}

// HTTP get request for a given Binance API endpoint
func (c *Client) request(endpoint string, parameters map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseUrl, endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("error while executing request: %v", err)
	}

	q := req.URL.Query()

	queryString := c.buildQueryString(q, parameters)
	if queryString != nil {
		req.URL.RawQuery = queryString.Encode()
	}

	req.Header.Add("X-MBX-APIKEY", c.APIKey)

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
