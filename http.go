package runhttp

import (
	"context"
	"net/http"
)

// HTTPConfig is the container for HTTP server configuration settings.
type HTTPConfig struct {
	Address string `description:"The listening address of the server."`
}

// Name returns the configuration root as it would appear in a config file.
func (*HTTPConfig) Name() string {
	return "httpserver"
}

// Description returns the help information for the configuration root.
func (*HTTPConfig) Description() string {
	return "HTTP server configuration."
}

// HTTPComponent implements the settings.Component interface for the HTTP server.
type HTTPComponent struct{}

// Settings returns a configuration with all defaults set.
func (*HTTPComponent) Settings() *HTTPConfig {
	return &HTTPConfig{
		Address: ":8080",
	}
}

// New produces a ServerFn bound to the given configuration.
func (*HTTPComponent) New(_ context.Context, conf *HTTPConfig) (ServerFn, error) {
	return func() *http.Server {
		return &http.Server{
			Addr: conf.Address,
		}
	}, nil
}
