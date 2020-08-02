package sesame

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	DummyAuthToken = "YOUR_AUTH_TOKEN"
)

func TestAuthorization(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if authHeader := req.Header.Get("Authorization"); authHeader != DummyAuthToken {
			http.Error(res, fmt.Sprintf("Invalid Authorization header: %s", authHeader), http.StatusUnauthorized)
			return
		}
		fmt.Fprint(res, "{}")
	}))
	defer sv.Close()

	cli := NewClient(DummyAuthToken)
	cli.BaseURL = sv.URL

	ctx := context.Background()

	s := Sesame{}
	err := cli.get(ctx, "test", nil, &s)
	assert.NoError(t, err)
}

func TestGetList(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var body []map[string]string
		el1 := map[string]string{
			"device_id": "00000000-0000-0000-0000-000000000000",
			"serial":    "ABC1234567",
			"nickname":  "Front door",
		}
		el2 := map[string]string{
			"device_id": "00000000-0000-0000-0000-000000000001",
			"serial":    "DEF7654321",
			"nickname":  "Back door",
		}
		body = append(body, el1, el2)
		p, _ := json.Marshal(body)
		fmt.Fprint(res, string(p))
	}))
	defer sv.Close()

	cli := NewClient(DummyAuthToken)
	cli.BaseURL = sv.URL

	ctx := context.Background()

	list, err := cli.GetList(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000000", list[0].DeviceID)
	assert.Equal(t, "ABC1234567", list[0].Serial)
	assert.Equal(t, "Front door", list[0].Nickname)
}

func TestGetStatus(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		var el map[string]interface{}
		s := strings.Split(req.RequestURI, "/")
		if s[2] == "device1" {
			el = map[string]interface{}{
				"locked":     true,
				"battery":    90,
				"responsive": true,
			}
		} else {
			el = map[string]interface{}{
				"locked":     false,
				"battery":    40,
				"responsive": true,
			}
		}
		p, _ := json.Marshal(el)
		fmt.Fprint(res, string(p))
	}))
	defer sv.Close()

	cli := NewClient(DummyAuthToken)
	cli.BaseURL = sv.URL

	ctx := context.Background()

	stat, err := cli.GetStatus(ctx, "device1")
	assert.NoError(t, err)
	assert.Equal(t, true, stat.Locked)
	assert.Equal(t, int64(90), stat.Battery)
	assert.Equal(t, true, stat.Responsive)

	stat, err = cli.GetStatus(ctx, "device2")
	assert.NoError(t, err)
	assert.Equal(t, false, stat.Locked)
	assert.Equal(t, int64(40), stat.Battery)
	assert.Equal(t, true, stat.Responsive)
}

func TestControl(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		el := map[string]string{
			"task_id": "01234567-890a-bcde-f012-34567890abcd",
		}
		p, _ := json.Marshal(el)
		fmt.Fprint(res, string(p))
	}))
	defer sv.Close()

	cli := NewClient(DummyAuthToken)
	cli.BaseURL = sv.URL

	ctx := context.Background()

	ctrl, err := cli.Control(ctx, "device1", "lock")
	assert.NoError(t, err)
	assert.Equal(t, "01234567-890a-bcde-f012-34567890abcd", ctrl.TaskID)
}

func TestGetExecutionResult(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		results := map[string]map[string]interface{}{
			"processing_task": {
				"status": "processing",
			},
			"terminated_task": {
				"status":     "terminated",
				"successful": true,
			},
			"error_task": {
				"status":     "terminated",
				"successful": false,
				"error":      "hogehoge",
			},
		}
		taskID := req.FormValue("task_id")
		p, _ := json.Marshal(results[taskID])
		fmt.Fprint(res, string(p))
	}))
	defer sv.Close()

	cli := NewClient(DummyAuthToken)
	cli.BaseURL = sv.URL

	ctx := context.Background()

	result, err := cli.GetExecutionResult(ctx, "processing_task")
	assert.NoError(t, err)
	assert.Equal(t, "processing", result.Status)

	result, err = cli.GetExecutionResult(ctx, "terminated_task")
	assert.NoError(t, err)
	assert.Equal(t, "terminated", result.Status)
	assert.Equal(t, true, result.Successful)

	result, err = cli.GetExecutionResult(ctx, "error_task")
	assert.NoError(t, err)
	assert.Equal(t, "terminated", result.Status)
	assert.Equal(t, false, result.Successful)
	assert.Equal(t, "hogehoge", result.Error)
}
