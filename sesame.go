package sesame

import (
	"context"
	"fmt"
	"net/url"
)

// Sesame represents sesame data.
type Sesame struct {
	DeviceID string `json:"device_id"`
	Serial   string `json:"serial"`
	Nickname string `json:"nickname"`
}

// SesameStatus represents status data.
type SesameStatus struct {
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
	GetStatus(ctx context.Context, deviceID string) (*SesameStatus, error)
	Control(ctx context.Context, deviceID string, command string) (*Control, error)
	GetExecutionResult(ctx context.Context, taskID string) (*ExecutionResult, error)
}

type apiClient struct {
	cli *Client
}

// NewAPIClient creates a SESAME SmartLock API client.
func NewAPIClient(accessToken string) *apiClient {
	cli := NewClient(accessToken)
	return &apiClient{cli: cli}
}

// GetList provides implementation of GET /sesames
// https://docs.candyhouse.co/#get-sesame-list
func (api *apiClient) GetList(ctx context.Context) ([]*Sesame, error) {
	var s []*Sesame
	if err := api.cli.Get(ctx, "sesames", nil, &s); err != nil {
		return nil, err
	}
	return s, nil
}

// GetStatus provides implementation of GET /sesame/{device_id}
// https://docs.candyhouse.co/#get-sesame-status
func (api *apiClient) GetStatus(ctx context.Context, deviceID string) (*SesameStatus, error) {
	var ss SesameStatus

	if deviceID == "" {
		return nil, fmt.Errorf("Invalid deviceID: %s", deviceID)
	}

	ep := fmt.Sprintf("sesame/%s", deviceID)
	if err := api.cli.Get(ctx, ep, nil, &ss); err != nil {
		return nil, err
	}
	return &ss, nil
}

// Control provides implementation of POST /sesame/{device_id}
// https://docs.candyhouse.co/#control-sesame
func (api *apiClient) Control(ctx context.Context, deviceID string, command string) (*Control, error) {
	var v Control
	if deviceID == "" {
		return nil, fmt.Errorf("Invalid deviceID: %s", deviceID)
	}

	cmd := ControlCommand{Command: command}

	ep := fmt.Sprintf("sesame/%s", deviceID)
	if err := api.cli.Post(ctx, ep, cmd, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

// GetExecutionResult provides implementation of GET /action-result?task_id={task_id}
// https://docs.candyhouse.co/#query-execution-result
func (api *apiClient) GetExecutionResult(ctx context.Context, taskID string) (*ExecutionResult, error) {
	var v ExecutionResult
	if taskID == "" {
		return nil, fmt.Errorf("Invalid taskID: %s", taskID)
	}

	p := url.Values{}
	p.Set("task_id", taskID)

	if err := api.cli.Get(ctx, "action-result", p, &v); err != nil {
		return nil, err
	}
	return &v, nil
}
