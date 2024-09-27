package runhttp

import (
	"context"
	"net/http"

	"github.com/asecurityteam/settings/v2"
)

// New uses the given source to generate a configured Runtime instace.
func New(ctx context.Context, s settings.Source, h http.Handler) (*Runtime, error) {
	runnerDst := new(Runtime)
	rt := NewComponent().WithHandler(h)
	err := settings.NewComponent(ctx, s, rt, runnerDst)
	return runnerDst, err
}
