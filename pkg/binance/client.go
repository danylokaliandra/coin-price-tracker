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
		return nil, fmt.Errorf("could not request time endpoint: %s", err.Error())
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body: %s", err.Error())
	}

	serverTime := ServerTime{}
	if err := json.Unmarshal(body, &serverTime); err != nil {
		return nil, fmt.Errorf("could not unmarshal server time: %s", err.Error())
	}

	return &serverTime, nil
}

func (c *Client) generateSignature(queryString string) (string, error) {
	mac := hmac.New(sha256.New, []byte(c.SecretKey))
	_, err := mac.Write([]byte(queryString))
	digest := mac.Sum(nil)
	if err != nil {
		return "", fmt.Errorf("could not generate hmac: %s", err.Error())
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
		return nil, fmt.Errorf("error while executing request: %s", err.Error())
	}

	q := req.URL.Query()
	queryString := c.buildQueryString(q, parameters)
	if queryString != nil {
		req.URL.RawQuery = queryString.Encode()
	}

	req.Header.Add("X-MBX-APIKEY", c.APIKey)

	response, err := c.HTTPClient.Do(req)
	if err != nil || response.StatusCode != 200 {
		return nil, fmt.Errorf("error while executing request: %s", err.Error())
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error while reading body: %s", err.Error())
	}

	return body, nil
}

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
			log.Warn("max retries exceeded...")
			break
		}

		retries += 1
		log.Warn("retrying... [%s/%s]: %s", retries, c.MaxRetries, err.Error())
		time.Sleep(time.Second * c.RetryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("could not request '%+v' endpoint: %s", endpoint, err.Error())
	}

	return body, err
}
