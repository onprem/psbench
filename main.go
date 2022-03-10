package main

import (
	"fmt"
	stdlog "log"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/onprem/psbench/pkg/bench"
	"github.com/onprem/psbench/pkg/config"
)

func run() error {
	cfg, err := config.ParseFlags()
	if err != nil {
		return fmt.Errorf("parsing config: %w", err)
	}

	logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	if cfg.LogFormat == "json" {
		logger = log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	}

	logger = level.NewFilter(logger, cfg.LogLevel)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	defer level.Info(logger).Log("msg", "exiting")

	_, err = bench.BenchPromscale(logger, cfg.Queries, cfg.Workers)
	if err != nil {
		return fmt.Errorf("benchmarking promscale: %w", err)
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		stdlog.Fatal(err)
	}
}
