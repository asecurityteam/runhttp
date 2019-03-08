package runhttp

import (
	"context"
	"net/http"
)

// Config is the top-level configuration container for
// a runtime.
type Config struct {
	HTTP      *HTTPConfig
	ConnState *ConnStateConfig
	Logger    *LoggerConfig
	Stats     *StatsConfig
	Signal    *SignalConfig
}

// Name returns the configuration root as it would appear in a config file.
func (*Config) Name() string {
	return "runtime"
}

// Component implements the settings.Component interface for an HTTP runtime.
type Component struct {
	Handler http.Handler
}

// Settings generates a configuration object with all defaults set.
func (*Component) Settings() *Config {
	return &Config{
		HTTP:      (&HTTPComponent{}).Settings(),
		ConnState: (&ConnStateComponent{}).Settings(),
		Logger:    (&LoggerComponent{}).Settings(),
		Stats:     (&StatsComponent{}).Settings(),
		Signal:    (&SignalComponent{}).Settings(),
	}
}

// New produces a configured runtime.
func (c *Component) New(ctx context.Context, conf *Config) (*Runtime, error) {
	log := &LoggerComponent{}
	connState := &ConnStateComponent{}
	stat := &StatsComponent{}
	sigs := &SignalComponent{}
	srv := &HTTPComponent{}

	logger, err := log.New(ctx, conf.Logger)
	if err != nil {
		return nil, err
	}
	stats, err := stat.New(ctx, conf.Stats)
	if err != nil {
		return nil, err
	}
	cs, err := connState.New(ctx, conf.ConnState)
	if err != nil {
		return nil, err
	}
	exit, err := sigs.New(ctx, conf.Signal)
	if err != nil {
		return nil, err
	}
	server, err := srv.New(ctx, conf.HTTP)
	if err != nil {
		return nil, err
	}

	return &Runtime{
		Logger:    logger,
		Stats:     stats,
		ConnState: cs,
		Exit:      exit,
		Server:    server,
		Handler:   c.Handler,
	}, nil
}
