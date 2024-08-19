package runhttp

import (
	"github.com/go-chi/chi/v5"
)

// RouterConfig is used as a simple default for NewDefaultRouter
type RouterConfig struct {
}

// NewDefaultRouter generates a mux.
// This version returns a mux from the chi project
// as a convenience for cases where custom middleware or additional
// routes need to be configured.
func NewDefaultRouter(conf *RouterConfig) *chi.Mux {
	router := chi.NewMux()
	healthCheckHandler := &HealthCheckHandler{}

	router.Get("/healthcheck", healthCheckHandler.Handle)

	return router
}
