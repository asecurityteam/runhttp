package runhttp

import (
	"context"
	"net/http"

	"github.com/rs/xstats"

	"github.com/asecurityteam/logevent"
)

// Logger is the project logging client interface. It is
// currently an alias to the logevent project.
type Logger = logevent.Logger

// LogFn is the type that should be accepted by components that
// intend to log content using the context logger.
type LogFn func(context.Context) Logger

// LoggerFromContext is the concrete implementation of LogFn
// that should be used at runtime.
var LoggerFromContext = logevent.FromContext

// Stat is the project metrics client interface. it is currently
// an alias for xstats.XStater.
type Stat = xstats.XStater

// StatFn is the type that should be accepted by components that
// intend to emit custom metrics using the context stat client.
type StatFn func(context.Context) Stat

// StatFromContext is the concrete implementation of StatFn that
// should be used at runtime.
var StatFromContext = xstats.FromContext

// SignalFn is a constructor for a shutdown signal. Signals are
// expected to be channels that emit an error that may be nil
// if the signal is expected.
type SignalFn func() chan error

// ServerFn is a constructor for *http.Server instances that will
// be hosted by the runtime.
type ServerFn func() *http.Server
