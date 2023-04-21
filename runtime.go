package runhttp

import (
	"context"
	"net/http"
	"time"

	connstate "github.com/asecurityteam/component-shared/connstate"
	expvar "github.com/asecurityteam/component-shared/expvar"
	signals "github.com/asecurityteam/component-shared/signals"
	hlog "github.com/asecurityteam/logevent/http"
	"github.com/rs/xstats"
)

// Runtime is the container for a restarting HTTP service. It will
// inject the given Logger and Stat into each request context, exit
// the Run method on each signal received from the output of the
// SignalFn, and use the ServerFn to regenerate a working server on
// subsequent Run calls.
type Runtime struct {
	Logger    Logger
	Stats     Stat
	ConnState *connstate.ConnState
	Expvar    *expvar.Expvar
	Exit      signals.Signal
	Server    *http.Server
	Handler   http.Handler
}

// Run the server until a signal is received.
func (r *Runtime) Run() error {

	go r.Expvar.Report()
	defer r.Expvar.Close()
	go r.ConnState.Report()
	defer r.ConnState.Close()

	handler := r.Handler
	handler = xstats.NewHandler(r.Stats, nil)(handler)
	handler = hlog.NewMiddleware(r.Logger)(handler)
	r.Server.Handler = handler

	go func() {
		r.Exit <- r.Server.ListenAndServe()
	}()

	err := <-r.Exit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = r.Server.Shutdown(ctx)

	return err
}
