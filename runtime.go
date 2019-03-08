package runhttp

import (
	"context"
	"net/http"
	"time"

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
	ConnState func() *ConnState
	Exit      SignalFn
	Server    ServerFn
	Handler   http.Handler
}

// Run the server until a signal is received.
func (r *Runtime) Run() error {
	exit := r.Exit()
	server := r.Server()
	cs := r.ConnState()
	cs.Stat = xstats.Copy(r.Stats)
	server.ConnState = cs.HandleEvent
	go cs.Report()
	handler := r.Handler
	handler = xstats.NewHandler(r.Stats, nil)(handler)
	handler = hlog.NewMiddleware(r.Logger)(handler)
	server.Handler = handler

	go func() {
		exit <- server.ListenAndServe()
	}()

	err := <-exit

	cs.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)

	return err
}
