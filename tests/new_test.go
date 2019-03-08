//+build integration

package tests

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/asecurityteam/settings"

	"github.com/asecurityteam/runhttp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type testLog struct {
	Message string `logevent:"message,default=test"`
}

func TestNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewMockLogger(ctrl)
	stat := NewMockStat(ctrl)
	ctx := context.Background()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := runhttp.LoggerFromContext(r.Context())
		logger.Info(testLog{})
		stat := runhttp.StatFromContext(r.Context())
		stat.Count("test", 1)
	})
	// Rather than mock out the settings.Source, it ends up being easier
	// to manage and slightly more realistic to use the ENV source but
	// populated with a static ENV list. This is easier because we don't
	// need to mock out the internal call structure of the settings.Source
	// which is largely irrelevant to this test. This is more realistic
	// because it leverages the public configuration API of the project
	// rather than internal knowledge of the settings project. For example,
	// these ENV vars are exactly the ones that users would set when running
	// the system.
	source, err := settings.NewEnvSource([]string{
		"RUNTIME_HTTPSERVER_ADDRESS=localhost:9090",
		"RUNTIME_LOGGER_OUTPUT=NULL",
		"RUNTIME_STATS_OUTPUT=NULL",
	})
	require.Nil(t, err)
	rt, err := runhttp.New(ctx, source, handler)
	require.Nil(t, err)
	// Swapping the NULL stat and logger so we can verify that the user
	// choice is propagated into the handler by leveraging the mock
	// expectations.
	rt.Logger = logger
	rt.Stats = stat

	// Ensure that the handler calls to the stat and logger
	// are triggered at-least one time to verify that our choice
	// was used when installing artifacts into the request context.
	// Using at-least once because the handler may be activated multiple
	// times as part of the test.
	logger.EXPECT().Copy().Return(logger).MinTimes(1)
	logger.EXPECT().Info(gomock.Any()).MinTimes(1)
	stat.EXPECT().Count("test", float64(1)).MinTimes(1)

	exit := make(chan error)
	go func() {
		exit <- rt.Run()
	}()

	// Ping the server until it is available or until we exceed a timeout
	// value. This is to account for arbitrary start-up time of the server
	// in the background.
	stop := time.Now().Add(5 * time.Second)
	for time.Now().Before(stop) {
		time.Sleep(100 * time.Millisecond)
		resp, err := http.DefaultClient.Get("http://localhost:9090")
		if err != nil {
			t.Log(err.Error())
			continue
		}
		if resp.StatusCode != http.StatusOK {
			t.Log(resp.StatusCode)
			continue
		}
		break
	}

	// The runtime establishes a signal handler for the entire
	// process. This means we have the process signal itself and
	// the runtime will intercept the call. This enables us to test
	// the signal based shutdown behavior.
	proc, _ := os.FindProcess(os.Getpid())
	proc.Signal(os.Interrupt)
	select {
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for exit")
	case err := <-exit:
		require.Nil(t, err)
	}
}
