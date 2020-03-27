package runhttp

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// RouterConfig is used as a simple default for NewDefaultRouter
type RouterConfig struct {
	// HealthCheck defines the route on which the service will respond
	// with automatic 200s. This is here to integrate with systems that
	// poll for liveliness. The default value is /healthcheck
	HealthCheck string
}

func applyDefaults(conf *RouterConfig) *RouterConfig {
	if conf.HealthCheck == "" {
		conf.HealthCheck = "/healthcheck"
	}
	return conf
}

// NewDefaultRouter generates a mux that already has AWS Lambda API
// routes bound. This version returns a mux from the chi project
// as a convenience for cases where custom middleware or additional
// routes need to be configured.
func NewDefaultRouter(conf *RouterConfig) *chi.Mux {
	conf = applyDefaults(conf)
	router := chi.NewMux()

	router.Use(middleware.Heartbeat(conf.HealthCheck))

	return router
}
