package runhttp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// MultiSignal converts a series of SignalFn instances into a single SignalFn.
func MultiSignal(fns []SignalFn) SignalFn {
	return func() chan error {
		c := make(chan error, len(fns))
		for _, fn := range fns {
			go func(fn SignalFn) {
				c <- <-fn()
			}(fn)
		}
		return c
	}
}

// OSSignal listens for OS signals SIGINT and SIGTERM and writes to the channel if
// either of those scenarios are encountered, nil is sent over the channel.
func OSSignal(signals []os.Signal) func() chan error {
	return func() chan error {
		c := make(chan os.Signal, len(signals))
		exit := make(chan error)
		go func() {
			signal.Notify(c, signals...)
			defer signal.Stop(c)
			<-c
			exit <- nil
		}()
		return exit
	}
}

// OSSignalConfig contains configuration for creating an OSSignal listener.
type OSSignalConfig struct {
	Signals []int `description:"Which signals to listen for."`
}

// Name of the configuration as it might appear in a file.
func (*OSSignalConfig) Name() string {
	return "os"
}

// Description of the configuration for help output.
func (*OSSignalConfig) Description() string {
	return "OS signal handlers for system shutdown."
}

// OSSignalComponent enables creation of an OS signal handler.
type OSSignalComponent struct{}

// Settings generates a default configuration.
func (*OSSignalComponent) Settings() *OSSignalConfig {
	return &OSSignalConfig{
		Signals: []int{int(syscall.SIGTERM), int(syscall.SIGINT)},
	}
}

// New generates a new SignalFn using OS signals.
func (*OSSignalComponent) New(_ context.Context, conf *OSSignalConfig) (SignalFn, error) {
	sigs := make([]os.Signal, 0, len(conf.Signals))
	for _, sig := range conf.Signals {
		sigs = append(sigs, syscall.Signal(sig))
	}
	return OSSignal(sigs), nil
}

// SignalConfig contains all configuration for enabling various shut down signals.
type SignalConfig struct {
	Installed []string `description:"Which signal handlers are installed. Choices are OS."`
	OS        *OSSignalConfig
}

// Name of the configuration as it appears in a file.
func (*SignalConfig) Name() string {
	return "signals"
}

// Description of the configuration for help output.
func (*SignalConfig) Description() string {
	return "Shutdown signal configuration."
}

// SignalComponent enables creation of signal handlers to shut down a system.
type SignalComponent struct{}

// Settings generates a default configuration.
func (*SignalComponent) Settings() *SignalConfig {
	return &SignalConfig{
		Installed: []string{"OS"},
		OS:        (&OSSignalComponent{}).Settings(),
	}
}

// New generates a new SignalFn using OS signals.
func (*SignalComponent) New(ctx context.Context, conf *SignalConfig) (SignalFn, error) {
	sigs := make([]SignalFn, 0, len(conf.Installed))
	for _, installed := range conf.Installed {
		switch installed {
		case "OS":
			sig, err := (&OSSignalComponent{}).New(ctx, conf.OS)
			if err != nil {
				return nil, err
			}
			sigs = append(sigs, sig)
		default:
			return nil, fmt.Errorf("unknown installed signal type %s", installed)
		}
	}
	return MultiSignal(sigs), nil
}
