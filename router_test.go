package runhttp

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouterHasHealthCheck(t *testing.T) {

	conf := &RouterConfig{}
	router := NewDefaultRouter(conf)

	resp := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "http://localhost/healthcheck", http.NoBody)
	router.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)
}
