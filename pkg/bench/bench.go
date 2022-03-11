package bench

import (
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/onprem/psbench/pkg/config"
)

type Result struct {
	NumQueries          uint
	TotalProcessingTime time.Duration
	MinQueryTime        time.Duration
	MaxQueryTime        time.Duration
	MedianQueryTime     time.Duration
	AverageQueryTime    time.Duration
}

func BenchPromscale(logger log.Logger, queries []config.Query, workers int) (Result, error) {
	queue := make(chan config.Query)

	// Send all queries to the task queue.
	go func() {
		for _, v := range queries {
			queue <- v
		}

		close(queue)
	}()

	var wg sync.WaitGroup

	// Run the specified number of workers to process queries for the queue.
	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			for q := range queue {
				time.Sleep(time.Second * 10)
				level.Info(logger).Log("msg", "executed query", "expression", q.Expression)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return Result{}, nil
}
