package binance

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Client struct {
	HTTPClient      *http.Client
	BaseUrl         string
	APIKey          string
	SecretKey       string
	RetryDelay      time.Duration
	MaxRetries      int
	TotalDaysOfData time.Duration
}

func NewClient(baseUrl string, apiKey string, secretKey string, timeout time.Duration, retryDelay time.Duration, maxRetries int, totalDaysOfData time.Duration) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * timeout,
		},
		BaseUrl:         baseUrl,
		APIKey:          apiKey,
		SecretKey:       secretKey,
		RetryDelay:      retryDelay,
		MaxRetries:      maxRetries,
		TotalDaysOfData: totalDaysOfData,
	}
}

func buildQueryString(params map[string]string) strings.Reader {
	values := url.Values{}
	for elem := range params {
		values.Add(elem, params[elem])
	}
	query := values.Encode()

	return *strings.NewReader(query)
}

func (c *Client) request(endpoint string) ([]byte, error) {
	parameters := buildQueryString(map[string]string{})

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseUrl, endpoint), &parameters)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-MBX-APIKEY", c.APIKey)

	response, err := c.HTTPClient.Do(req)
	if err != nil || response.StatusCode != 200 {
		return nil, fmt.Errorf("error while executing request")
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *Client) RequestWithRetries(endpoint string) ([]byte, error) {
	var body []byte
	var err error
	retries := 0

	for {
		body, err = c.request(endpoint)
		if err == nil || retries == c.MaxRetries {
			break
		}

		retries += 1
		fmt.Fprintf(os.Stderr, "Request error: %+v\n", err)
		fmt.Fprintf(os.Stderr, "Retrying in: %+v\n", c.RetryDelay)
		fmt.Fprintf(os.Stderr, "Retry : %+v/%+v\n", retries, c.MaxRetries)
		time.Sleep(time.Second * c.RetryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("could not request '%+v' endpoint", endpoint)
	}

	return body, err
}
