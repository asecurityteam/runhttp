package runhttp

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	statCounterClientNew      = "http.server.connstate.new"
	statGaugeClientNew        = "http.server.connstate.new.gauge"
	statCounterClientActive   = "http.server.connstate.active"
	statGaugeClientActive     = "http.server.connstate.active.gauge"
	statCounterClientIdle     = "http.server.connstate.idle"
	statGaugeClientIdle       = "http.server.connstate.idle.gauge"
	statCounterClientClosed   = "http.server.connstate.closed"
	statCounterClientHijacked = "http.server.connstate.hijacked"
	connStateInterval         = 5 * time.Second
)

// ConnState plugs into the http.Server.ConnState attribute to track the number of
// client connections to the server.
type ConnState struct {
	Stat                      Stat
	Tracking                  *sync.Map
	NewClientCounterName      string
	NewClientGaugeName        string
	ActiveClientCounterName   string
	ActiveClientGaugeName     string
	IdleClientCounterName     string
	IdleClientGaugeName       string
	ClosedClientCounterName   string
	HijackedClientCounterName string
	Interval                  time.Duration
	statMut                   *sync.Mutex
	stopMut                   *sync.Mutex
	stop                      bool
}

// Report loops on a time interval and pushes a set of gauge metrics.
func (c *ConnState) Report() {
	ticker := time.NewTicker(c.Interval)
	defer ticker.Stop()
	for range ticker.C {
		c.report()
		c.stopMut.Lock()
		if c.stop {
			c.stopMut.Unlock()
			return
		}
		c.stopMut.Unlock()
	}
}

// Close the reporting loop.
func (c *ConnState) Close() {
	c.stopMut.Lock()
	defer c.stopMut.Unlock()
	c.stop = true
}

func (c *ConnState) report() {
	var n float64
	var a float64
	var i float64
	c.Tracking.Range(func(key interface{}, value interface{}) bool {
		switch value.(http.ConnState) {
		case http.StateNew:
			n = n + 1
		case http.StateActive:
			a = a + 1
		case http.StateIdle:
			i = i + 1
		}
		return true
	})
	c.statMut.Lock()
	defer c.statMut.Unlock()
	c.Stat.Gauge(c.NewClientGaugeName, n)
	c.Stat.Gauge(c.ActiveClientGaugeName, a)
	c.Stat.Gauge(c.IdleClientGaugeName, i)
}

// HandleEvent tracks state changes of a connection.
func (c *ConnState) HandleEvent(conn net.Conn, state http.ConnState) {
	c.statMut.Lock()
	defer c.statMut.Unlock()
	switch state {
	case http.StateNew:
		c.Stat.Count(c.NewClientCounterName, 1)
		c.Tracking.Store(conn, state)
	case http.StateActive:
		c.Stat.Count(c.ActiveClientCounterName, 1)
		c.Tracking.Store(conn, state)
	case http.StateIdle:
		c.Stat.Count(c.IdleClientCounterName, 1)
		c.Tracking.Store(conn, state)
	case http.StateHijacked:
		c.Stat.Count(c.HijackedClientCounterName, 1)
		c.Tracking.Delete(conn)
	case http.StateClosed:
		c.Stat.Count(c.ClosedClientCounterName, 1)
		c.Tracking.Delete(conn)
	}
}

// ConnStateConfig is a container for internal metrics settings.
type ConnStateConfig struct {
	NewCounter      string        `description:"Name of the counter metric tracking new clients."`
	NewGauge        string        `description:"Name of the gauge metric tracking new clients."`
	ActiveCounter   string        `description:"Name of the counter metric tracking active clients."`
	ActiveGauge     string        `description:"Name of the gauge metric tracking active clients."`
	IdleCounter     string        `description:"Name of the counter metric tracking idle clients."`
	IdleGauge       string        `description:"Name of the gauge metric tracking idle clients."`
	ClosedCounter   string        `description:"Name of the counter metric tracking closed clients."`
	HijackedCounter string        `description:"Name of the counter metric tracking hijacked clients."`
	ReportInterval  time.Duration `description:"Interval on which gauges are reported."`
}

// Name of the configuration root.
func (*ConnStateConfig) Name() string {
	return "connstate"
}

// Description returns the help information for the configuration root.
func (*ConnStateConfig) Description() string {
	return "Connection state metric names."
}

// ConnStateComponent implements the settings.Component interface for connection
// state monitoring.
type ConnStateComponent struct{}

// Settings returns a configuration with all defaults set.
func (*ConnStateComponent) Settings() *ConnStateConfig {
	return &ConnStateConfig{
		NewCounter:      statCounterClientNew,
		NewGauge:        statGaugeClientNew,
		ActiveCounter:   statCounterClientActive,
		ActiveGauge:     statGaugeClientActive,
		IdleCounter:     statCounterClientIdle,
		IdleGauge:       statGaugeClientIdle,
		ClosedCounter:   statCounterClientClosed,
		HijackedCounter: statCounterClientHijacked,
		ReportInterval:  connStateInterval,
	}
}

// New produces a ServerFn bound to the given configuration.
func (*ConnStateComponent) New(_ context.Context, conf *ConnStateConfig) (func() *ConnState, error) {
	return func() *ConnState {
		return &ConnState{
			Tracking:                  &sync.Map{},
			NewClientCounterName:      conf.NewCounter,
			NewClientGaugeName:        conf.NewGauge,
			ActiveClientCounterName:   conf.ActiveCounter,
			ActiveClientGaugeName:     conf.ActiveGauge,
			IdleClientCounterName:     conf.IdleCounter,
			IdleClientGaugeName:       conf.IdleGauge,
			ClosedClientCounterName:   conf.ClosedCounter,
			HijackedClientCounterName: conf.HijackedCounter,
			Interval:                  conf.ReportInterval,
			statMut:                   &sync.Mutex{},
			stopMut:                   &sync.Mutex{},
		}
	}, nil

}
