// Handles HTTP requests logic
package coingecko

import (
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
	RetryDelay time.Duration // in seconds
	MaxRetries int
}

// Create a new Client object
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * viper.GetDuration("tracklet.timeout"),
		},
		BaseUrl:    viper.GetString("aggregators.coingecko.apiBaseUrl"),
		RetryDelay: viper.GetDuration("tracklet.retryDelay"),
		MaxRetries: viper.GetInt("tracklet.maxRetries"),
	}
}

// Build http query string
func (c *Client) buildQueryString(q url.Values, params map[string]string) url.Values {
	if len(params) == 0 {
		return nil
	}
	for elem := range params {
		q.Add(elem, params[elem])
	}

	return q
}

// HTTP get request for a given CoinGecko API endpoint
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
