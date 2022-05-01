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
	"os"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	BaseUrl    string
	APIKey     string
	SecretKey  string
	RetryDelay time.Duration
	MaxRetries int
}

type ServerTime struct {
	ServerTime int64 `json:"serverTime"`
}

func NewClient(baseUrl string, apiKey string, secretKey string, timeout time.Duration, retryDelay time.Duration, maxRetries int) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: time.Second * timeout,
		},
		BaseUrl:    baseUrl,
		APIKey:     apiKey,
		SecretKey:  secretKey,
		RetryDelay: retryDelay,
		MaxRetries: maxRetries,
	}
}

func (c *Client) getServerTime() (*ServerTime, error) {
	response, err := http.Get(fmt.Sprintf("%s/api/v3/time", c.BaseUrl))
	if err != nil {
		return nil, fmt.Errorf("could not request time endpoint")
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body")
	}

	serverTime := ServerTime{}
	if err := json.Unmarshal(body, &serverTime); err != nil {
		return nil, fmt.Errorf("could not unmarshal server time %+v", err)
	}

	return &serverTime, nil
}

func (c *Client) generateSignature(queryString string) (string, error) {
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	_, err := mac.Write([]byte(queryString))
	digest := mac.Sum(nil)
	if err != nil {
		return "", fmt.Errorf("could not generate hmac")
	}
	hexDigest := hex.EncodeToString(digest)
	return hexDigest, nil
}

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

func (c *Client) request(endpoint string, parameters map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.BaseUrl, endpoint), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	queryString := c.buildQueryString(q, parameters)
	if queryString != nil {
		req.URL.RawQuery = queryString.Encode()
	}

	req.Header.Add("X-MBX-APIKEY", c.APIKey)

	response, err := c.HTTPClient.Do(req)
	if err != nil || response.StatusCode != 200 {
		return nil, fmt.Errorf("error while executing request: %+v", response.Status)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body")
	}

	return body, nil
}

func (c *Client) RequestWithRetries(endpoint string, parameters map[string]string) ([]byte, error) {
	var body []byte
	var err error
	retries := 0

	for {
		body, err = c.request(endpoint, parameters)
		if err == nil || retries == c.MaxRetries {
			break
		}

		retries += 1
		fmt.Fprintf(os.Stderr, "%+v (Retry %+v/%+v)\n", err, retries, c.MaxRetries)
		time.Sleep(time.Second * c.RetryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("could not request '%+v' endpoint", endpoint)
	}

	return body, err
}
