package runhttp

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/asecurityteam/logevent"
)

const (
	defaultLevel  = "INFO"
	defaultOutput = "STDOUT"
)

// LoggerConfig contains all configuration values for creating
// a system logger.
type LoggerConfig struct {
	Level  string `description:"The minimum level of logs to emit. One of DEBUG, INFO, WARN, ERROR."`
	Output string `description:"Destination stream of the logs. One of STDOUT, NULL."`
}

// Name of the configuration as it might appear in config files.
func (*LoggerConfig) Name() string {
	return "logger"
}

// LoggerComponent enables creating configured loggers.
type LoggerComponent struct{}

// Settings generates a LoggerConfig with default values applied.
func (*LoggerComponent) Settings() *LoggerConfig {
	return &LoggerConfig{
		Level:  defaultLevel,
		Output: defaultOutput,
	}
}

// New creates a configured logger instance.
func (*LoggerComponent) New(_ context.Context, conf *LoggerConfig) (Logger, error) {
	var output io.Writer
	switch conf.Output {
	case "STDOUT":
		output = os.Stdout
	case "NULL":
		output = ioutil.Discard
	default:
		return nil, fmt.Errorf("unknown logger output %s", conf.Output)
	}
	return logevent.New(logevent.Config{Level: conf.Level, Output: output}), nil
}
