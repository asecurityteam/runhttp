package runhttp

import (
	"context"
	"net/http"

	connstate "github.com/asecurityteam/component-shared/connstate"
	expvar "github.com/asecurityteam/component-shared/expvar"
	log "github.com/asecurityteam/component-shared/log"
	signals "github.com/asecurityteam/component-shared/signals"
	stat "github.com/asecurityteam/component-shared/stat"
	"github.com/rs/xstats"
)

// Config is the top-level configuration container for
// a runtime.
type Config struct {
	HTTP      *HTTPConfig
	ConnState *connstate.Config
	Expvar    *expvar.Config
	Logger    *log.Config
	Stats     *stat.Config
	Signal    *signals.Config
}

// Name returns the configuration root as it would appear in a config file.
func (*Config) Name() string {
	return "runtime"
}

// Component implements the settings.Component interface for an HTTP runtime.
type Component struct {
	HTTP      *HTTPComponent
	Connstate *connstate.Component
	Expvar    *expvar.Component
	Logger    *log.Component
	Stats     *stat.Component
	Signal    *signals.Component
	Handler   http.Handler
}

// NewComponent populates the component with some default values.
func NewComponent() *Component {
	return &Component{
		HTTP:      &HTTPComponent{},
		Connstate: connstate.NewComponent(),
		Expvar:    expvar.NewComponent(),
		Logger:    log.NewComponent(),
		Stats:     stat.NewComponent(),
		Signal:    signals.NewComponent(),
	}
}

// WithHandler returns a copy of the component bound to the given handler.
func (c *Component) WithHandler(h http.Handler) *Component {
	n := NewComponent()
	n.Handler = h
	return n
}

// Settings generates a configuration object with all defaults set.
func (c *Component) Settings() *Config {
	return &Config{
		HTTP:      c.HTTP.Settings(),
		ConnState: c.Connstate.Settings(),
		Expvar:    c.Expvar.Settings(),
		Logger:    c.Logger.Settings(),
		Stats:     c.Stats.Settings(),
		Signal:    c.Signal.Settings(),
	}
}

// New produces a configured runtime.
func (c *Component) New(ctx context.Context, conf *Config) (*Runtime, error) {
	logger, err := c.Logger.New(ctx, conf.Logger)
	if err != nil {
		return nil, err
	}
	stats, err := c.Stats.New(ctx, conf.Stats)
	if err != nil {
		return nil, err
	}
	cs, err := c.Connstate.WithStat(xstats.Copy(stats)).New(ctx, conf.ConnState)
	if err != nil {
		return nil, err
	}
	expvar, err := c.Expvar.WithStat(xstats.Copy(stats)).New(ctx, conf.Expvar)
	if err != nil {
		return nil, err
	}
	exit, err := c.Signal.New(ctx, conf.Signal)
	if err != nil {
		return nil, err
	}
	server, err := c.HTTP.New(ctx, conf.HTTP)
	if err != nil {
		return nil, err
	}
	server.ConnState = cs.HandleEvent

	return &Runtime{
		Logger:    logger,
		Stats:     stats,
		ConnState: cs,
		Expvar:    expvar,
		Exit:      exit,
		Server:    server,
		Handler:   c.Handler,
	}, nil
}
