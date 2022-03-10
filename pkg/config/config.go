package config

import (
	"flag"
	"fmt"

	"github.com/go-kit/kit/log/level"
)

type Config struct {
	LogFormat string
	LogLevel  level.Option
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}

	// Logger flags.
	logLevelRaw := flag.String("log.level", "info", "The log filtering level. Options: 'error', 'warn', 'info', 'debug'.")
	flag.StringVar(&cfg.LogFormat, "log.format", "logfmt", "The log format to use. Options: 'logfmt', 'json'.")

	flag.Parse()

	ll, err := parseLogLevel(logLevelRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	cfg.LogLevel = ll

	return cfg, nil
}

func parseLogLevel(logLevelRaw *string) (level.Option, error) {
	switch *logLevelRaw {
	case "error":
		return level.AllowError(), nil
	case "warn":
		return level.AllowWarn(), nil
	case "info":
		return level.AllowInfo(), nil
	case "debug":
		return level.AllowDebug(), nil
	default:
		return nil, fmt.Errorf("unexpected log level: %s", *logLevelRaw) //nolint:goerr113
	}
}
