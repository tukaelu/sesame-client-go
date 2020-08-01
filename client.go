package sesame

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unsafe"
)

const (
	apiBaseURL = "https://api.candyhouse.co/public"
	libVersion = "0.1.0"
	reqTimeout = 30 * time.Second
)

// Client for SESAME SmartLock API.
type Client struct {
	BaseURL     string
	AccessToken string
	UserAgent   string
	HTTPClient  *http.Client
}

// NewClient creates a client for the SESAME SmartLock API
func NewClient(accessToken string) *Client {
	return &Client{
		BaseURL:     apiBaseURL,
		AccessToken: accessToken,
		UserAgent:   fmt.Sprintf("tukaelu/sesame-client-go (Ver: %s)", libVersion),
		HTTPClient:  &http.Client{},
	}
}

// Get is an implementation of the HTTP GET method.
func (cli *Client) Get(ctx context.Context, path string, params url.Values, p interface{}) error {
	endpoint := fmt.Sprintf("%s/%s", cli.BaseURL, path)
	if params != nil {
		endpoint = fmt.Sprintf("%s?%s", endpoint, params.Encode())
	}

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	res, err := cli.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if !(res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices) {
		reason, err := ioutil.ReadAll(res.Body)
		if err != nil || len(reason) == 0 {
			return fmt.Errorf("Request failed: Status=%d (no reason)", res.StatusCode)
		}
		return fmt.Errorf("Request failed: Status=%d, Error= %s", res.StatusCode, string(reason))
	}

	if err := json.NewDecoder(res.Body).Decode(p); err != nil {
		return fmt.Errorf("Failed to parse the response. (%s)", err.Error())
	}

	return nil
}

// Post is an implementation of the HTTP POST method.
func (cli *Client) Post(ctx context.Context, path string, params interface{}, p interface{}) error {
	endpoint := fmt.Sprintf("%s/%s", cli.BaseURL, path)

	j, err := json.Marshal(params)
	if err != nil {
		return err
	}
	payload := strings.NewReader(*(*string)(unsafe.Pointer(&j)))

	req, err := http.NewRequest(http.MethodPost, endpoint, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := cli.doRequest(ctx, req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if !(res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusMultipleChoices) {
		reason, err := ioutil.ReadAll(res.Body)
		if err != nil || len(reason) == 0 {
			return fmt.Errorf("Request failed: Status=%d (no reason)", res.StatusCode)
		}
		return fmt.Errorf("Request failed: Status=%d, Error= %s", res.StatusCode, string(reason))
	}

	if err := json.NewDecoder(res.Body).Decode(p); err != nil {
		return fmt.Errorf("Failed to parse the response. (%s)", err.Error())
	}

	return nil
}

func (cli *Client) doRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", fmt.Sprintf("%s", cli.AccessToken))
	req.Header.Set("User-Agent", cli.UserAgent)
	cli.HTTPClient.Timeout = reqTimeout
	return cli.HTTPClient.Do(req)
}
