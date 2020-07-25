package opensesame

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

	api := NewSesameAPI(DummyAuthToken)
	api.cli.BaseURL = sv.URL

	ctx := context.Background()

	s := Sesame{}
	err := api.cli.Get(ctx, "test", nil, &s)
	assert.NoError(t, err)
}
