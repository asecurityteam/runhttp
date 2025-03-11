package runhttp

import (
	"net/http"
)

// RouterConfig is used as a simple default for NewDefaultRouter
type RouterConfig struct {
}

// NewDefaultRouter generates a mux.
// This version returns a mux from the chi project
// as a convenience for cases where custom middleware or additional
// routes need to be configured.
func NewDefaultRouter(conf *RouterConfig) *http.ServeMux {
	router := http.NewServeMux()
	healthCheckHandler := &HealthCheckHandler{}

	router.HandleFunc("/healthcheck", healthCheckHandler.Handle)

	return router
}
