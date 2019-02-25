package runhttp

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/rs/xstats"
	"github.com/rs/xstats/dogstatsd"
)

const (
	defaultStatsOutput  = "DATADOG"
	defaultDDAddr       = "localhost:8125"
	defaultDDPacketSize = 1 << 15 // Matches what is used in the xstats library.
	defaultDDFlush      = 10 * time.Second
)

var (
	defaultDDTags = []string{}
)

// NullStatsConfig is empty. There are no options for NULL.
type NullStatsConfig struct{}

// Name of the configuration as it might appear in config files.
func (*NullStatsConfig) Name() string {
	return "nullstat"
}

// NullStatsComponent implements the settings.Component interface for
// a NOP stat client.
type NullStatsComponent struct{}

// Settings generates a config with default values applied.
func (*NullStatsComponent) Settings() *NullStatsConfig {
	return &NullStatsConfig{}
}

// New creates a configured stats client.
func (*NullStatsComponent) New(_ context.Context, conf *NullStatsConfig) (Stat, error) {
	return xstats.FromContext(context.Background()), nil
}

// DatadogStatsConfig is for configuration a datadog client.
type DatadogStatsConfig struct {
	Address       string        `description:"Listener address to use when sending metrics."`
	FlushInterval time.Duration `description:"Frequencing of sending metrics to listener."`
	Tags          []string      `description:"Any static tags for all metrics."`
	PacketSize    int           `description:"Max packet size to send."`
}

// Name of the configuration as it might appear in config files.
func (*DatadogStatsConfig) Name() string {
	return "datadog"
}

// DatadogStatsComponent implements the settings.Component interface for
// a datadog stats client.
type DatadogStatsComponent struct{}

// Settings generates a config with default values applied.
func (*DatadogStatsComponent) Settings() *DatadogStatsConfig {
	return &DatadogStatsConfig{
		Address:       defaultDDAddr,
		FlushInterval: defaultDDFlush,
		Tags:          defaultDDTags,
		PacketSize:    defaultDDPacketSize,
	}
}

// New creates a configured stats client.
func (*DatadogStatsComponent) New(_ context.Context, conf *DatadogStatsConfig) (Stat, error) {
	writer, err := net.Dial("udp", conf.Address)
	if err != nil {
		return nil, err
	}
	return xstats.New(dogstatsd.NewMaxPacket(writer, conf.FlushInterval, conf.PacketSize)), nil
}

// StatsConfig contains all configuration values for creating
// a system stat client.
type StatsConfig struct {
	Output   string `description:"Destination stream of the stats. One of NULLSTAT, DATADOG."`
	NullStat *NullStatsConfig
	Datadog  *DatadogStatsConfig
}

// Name of the configuration as it might appear in config files.
func (*StatsConfig) Name() string {
	return "stats"
}

// StatsComponent enables creating configured loggers.
type StatsComponent struct{}

// Settings generates a StatsConfig with default values applied.
func (*StatsComponent) Settings() *StatsConfig {
	return &StatsConfig{
		Output:   defaultStatsOutput,
		NullStat: (&NullStatsComponent{}).Settings(),
		Datadog:  (&DatadogStatsComponent{}).Settings(),
	}
}

// New creates a configured stats client.
func (*StatsComponent) New(ctx context.Context, conf *StatsConfig) (Stat, error) {
	switch conf.Output {
	case "NULL":
		return (&NullStatsComponent{}).New(ctx, conf.NullStat)
	case "DATADOG":
		return (&DatadogStatsComponent{}).New(ctx, conf.Datadog)
	default:
		return nil, fmt.Errorf("unknown stats output %s", conf.Output)
	}
}
