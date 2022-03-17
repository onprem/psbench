package bench

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"

	"github.com/onprem/psbench/pkg/config"
)

func BenchPromscale(logger log.Logger, address string, queries []config.Query, workers int) (Result, error) {
	queue := make(chan config.Query)

	// Send all queries to the task queue.
	go func() {
		for _, v := range queries {
			queue <- v
		}

		close(queue)
	}()

	c, err := api.NewClient(api.Config{Address: address})
	if err != nil {
		return Result{}, fmt.Errorf("error creating prometheus api client: %w", err)
	}

	promapi := v1.NewAPI(c)

	queryTimes := make([]time.Duration, 0, len(queries))
	timeQ := make(chan time.Duration)

	go func() {
		for d := range timeQ {
			queryTimes = sortedAppend(queryTimes, d)
		}
	}()

	failed := 0
	failedQ := make(chan struct{})

	go func() {
		for _ = range failedQ {
			failed += 1
		}
	}()

	var wg sync.WaitGroup

	// Run the specified number of workers to process queries for the queue.
	for i := 0; i < workers; i++ {
		wg.Add(1)

		go func() {
			for q := range queue {
				start := time.Now()

				_, _, err := promapi.QueryRange(
					context.TODO(),
					q.Expression,
					v1.Range{Start: q.StartTime, End: q.EndTime, Step: q.Step},
				)
				if err != nil {
					level.Error(logger).Log("msg", "executing query", "expression", q.Expression, "err", err)
					failedQ <- struct{}{}
					continue
				}

				timeQ <- time.Since(start)

				level.Debug(logger).Log("msg", "executed query", "expression", q.Expression)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	return generateResult(queryTimes, len(queries), failed), nil
}

// generateResult evaluates the given set of query times and generates the Result struct from it.
// It expects the queryTimes array to be already sorted.
func generateResult(queryTimes []time.Duration, total, failed int) Result {
	if len(queryTimes) == 0 {
		return Result{}
	}

	res := Result{
		NumQueries:        total,
		SuccessfulQueries: len(queryTimes),
		FailedQueries:     failed,
	}

	for _, v := range queryTimes {
		res.TotalProcessingTime += v
	}

	res.AverageQueryTime = res.TotalProcessingTime / time.Duration(res.SuccessfulQueries)

	res.MinQueryTime = queryTimes[0]
	res.MaxQueryTime = queryTimes[len(queryTimes)-1]

	mIdx := len(queryTimes) / 2

	res.MedianQueryTime = queryTimes[mIdx]
	if len(queryTimes)%2 == 0 {
		res.MedianQueryTime = (queryTimes[mIdx-1] + queryTimes[mIdx]) / 2
	}

	return res
}

func sortedAppend(arr []time.Duration, d time.Duration) []time.Duration {
	arr = append(arr, d)
	for i := len(arr) - 1; i > 0; i-- {
		if arr[i] < arr[i-1] {
			arr[i], arr[i-1] = arr[i-1], arr[i]
		} else {
			// Break the circuit once we reach the already sorted part of array.
			break
		}
	}
	return arr
}

type Result struct {
	NumQueries          int
	SuccessfulQueries   int
	FailedQueries       int
	TotalProcessingTime time.Duration
	MinQueryTime        time.Duration
	MaxQueryTime        time.Duration
	MedianQueryTime     time.Duration
	AverageQueryTime    time.Duration
}

func (r Result) String() string {
	return fmt.Sprintf(
		"Total Number of Queries: \t%d\n"+
			"Successful Queries: \t\t%d\n"+
			"Failed Queries: \t\t%d\n"+
			"Total Processing Time: \t\t%v\n"+
			"Minimum Query Time: \t\t%v\n"+
			"Maximum Query Time: \t\t%v\n"+
			"Median Query Time: \t\t%v\n"+
			"Avergae Query Time: \t\t%v",
		r.NumQueries, r.SuccessfulQueries, r.FailedQueries,
		r.TotalProcessingTime, r.MinQueryTime, r.MaxQueryTime, r.MedianQueryTime, r.AverageQueryTime,
	)
}
