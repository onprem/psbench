package config

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/go-kit/kit/log/level"
)

type Config struct {
	LogFormat string
	LogLevel  level.Option

	Workers int
	Queries []Query
}

type Query struct {
	Expression string
	StartTime  uint64
	EndTime    uint64
	Step       uint64
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}

	// Logger flags.
	logLevelRaw := flag.String("log.level", "info", "The log filtering level. Options: 'error', 'warn', 'info', 'debug'.")
	flag.StringVar(&cfg.LogFormat, "log.format", "logfmt", "The log format to use. Options: 'logfmt', 'json'.")

	flag.IntVar(&cfg.Workers, "workers", 1,
		"The number of workers/clients to run parallelly to query the Promscale instance.")

	// Queries file flags.
	queriesFilePath := flag.String(
		"queries.file",
		"",
		"Path to CSV file that contains rows with queries to execute. Format: PromQL query,start_time,end_time,step size.",
	)

	flag.Parse()

	ll, err := parseLogLevel(logLevelRaw)
	if err != nil {
		return nil, fmt.Errorf("invalid log level: %w", err)
	}

	cfg.LogLevel = ll

	if cfg.Workers < 1 {
		return nil, fmt.Errorf("number of workers needs to be greater than 0")
	}

	queriesFile, err := os.Open(*queriesFilePath)
	if err != nil {
		return nil, fmt.Errorf("opening queries file: %w", err)
	}
	defer queriesFile.Close()

	if cfg.Queries, err = parseQueriesFile(queriesFile); err != nil {
		return nil, fmt.Errorf("parsing queries file: %w", err)
	}

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

func parseQueriesFile(file io.Reader) ([]Query, error) {
	qr := csv.NewReader(file)
	qr.LazyQuotes = true
	qr.Comma = '|'

	queries := make([]Query, 0)

	for {
		r, err := qr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}

			return queries, err
		}

		if len(r) < 4 {
			return queries, fmt.Errorf("invalid number of rows in CSV")
		}

		q := Query{Expression: r[0]}

		if q.StartTime, err = strconv.ParseUint(r[1], 10, 64); err != nil {
			return queries, err
		}

		if q.EndTime, err = strconv.ParseUint(r[2], 10, 64); err != nil {
			return queries, err
		}

		if q.Step, err = strconv.ParseUint(r[3], 10, 64); err != nil {
			return queries, err
		}

		queries = append(queries, q)
	}
}
