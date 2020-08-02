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

// Sesame represents sesame data.
type Sesame struct {
	DeviceID string `json:"device_id"`
	Serial   string `json:"serial"`
	Nickname string `json:"nickname"`
}

// Status represents status data.
type Status struct {
	Locked     bool  `json:"locked"`
	Battery    int64 `json:"battery"`
	Responsive bool  `json:"responsive"`
}

// ControlCommand represents command data.
type ControlCommand struct {
	Command string `json:"command"`
}

// Control represents control task id.
type Control struct {
	TaskID string `json:"task_id"`
}

// ExecutionResult represents execution status.
type ExecutionResult struct {
	Status     string `json:"status"`
	Successful bool   `json:"successful"`
	Error      string `json:"error"`
}

// SesameAPI provides interface of Sesame API.
type SesameAPI interface {
	GetList(ctx context.Context) ([]*Sesame, error)
	GetStatus(ctx context.Context, deviceID string) (*Status, error)
	Control(ctx context.Context, deviceID string, command string) (*Control, error)
	GetExecutionResult(ctx context.Context, taskID string) (*ExecutionResult, error)
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

// GetList provides implementation of GET /sesames
// https://docs.candyhouse.co/#get-sesame-list
func (cli *Client) GetList(ctx context.Context) ([]*Sesame, error) {
	var s []*Sesame
	if err := cli.get(ctx, "sesames", nil, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetStatus provides implementation of GET /sesame/{device_id}
// https://docs.candyhouse.co/#get-sesame-status
func (cli *Client) GetStatus(ctx context.Context, deviceID string) (*Status, error) {
	var s Status

	if deviceID == "" {
		return nil, fmt.Errorf("Invalid deviceID: %s", deviceID)
	}

	ep := fmt.Sprintf("sesame/%s", deviceID)
	if err := cli.get(ctx, ep, nil, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

// Control provides implementation of POST /sesame/{device_id}
// https://docs.candyhouse.co/#control-sesame
func (cli *Client) Control(ctx context.Context, deviceID string, command string) (*Control, error) {
	var v Control
	if deviceID == "" {
		return nil, fmt.Errorf("Invalid deviceID: %s", deviceID)
	}

	cmd := ControlCommand{Command: command}

	ep := fmt.Sprintf("sesame/%s", deviceID)
	if err := cli.post(ctx, ep, cmd, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// GetExecutionResult provides implementation of GET /action-result?task_id={task_id}
// https://docs.candyhouse.co/#query-execution-result
func (cli *Client) GetExecutionResult(ctx context.Context, taskID string) (*ExecutionResult, error) {
	var v ExecutionResult
	if taskID == "" {
		return nil, fmt.Errorf("Invalid taskID: %s", taskID)
	}

	p := url.Values{}
	p.Set("task_id", taskID)

	if err := cli.get(ctx, "action-result", p, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// Get is an implementation of the HTTP GET method.
func (cli *Client) get(ctx context.Context, path string, params url.Values, p interface{}) error {
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
func (cli *Client) post(ctx context.Context, path string, params interface{}, p interface{}) error {
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
