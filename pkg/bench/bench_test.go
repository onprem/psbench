package bench

import (
	"testing"
	"time"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestGenerateResult(t *testing.T) {
	cases := []struct {
		name  string
		input []time.Duration
		exp   Result
	}{
		{
			name:  "empty input",
			input: []time.Duration{},
			exp:   Result{},
		},
		{
			name:  "sample case",
			input: []time.Duration{time.Second, 1 * time.Second, 2 * time.Second, 4 * time.Second, 9 * time.Second},
			exp: Result{
				SuccessfulQueries:   5,
				TotalProcessingTime: 17 * time.Second,
				MinQueryTime:        1 * time.Second,
				MaxQueryTime:        9 * time.Second,
				MedianQueryTime:     2 * time.Second,
				AverageQueryTime:    (17 * time.Second) / 5,
			},
		},
		{
			name:  "even observations",
			input: []time.Duration{time.Second, 1 * time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 9 * time.Second},
			exp: Result{
				SuccessfulQueries:   6,
				TotalProcessingTime: 20 * time.Second,
				MinQueryTime:        1 * time.Second,
				MaxQueryTime:        9 * time.Second,
				MedianQueryTime:     (5 * time.Second) / 2,
				AverageQueryTime:    (20 * time.Second) / 6,
			},
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			got := generateResult(v.input, 0, 0)

			testutil.Equals(t, v.exp, got)
		})
	}
}

func TestSortedAppend(t *testing.T) {
	cases := []struct {
		name   string
		arr    []time.Duration
		insert []time.Duration
		exp    []time.Duration
	}{
		{
			name:   "empty input, one insert",
			arr:    []time.Duration{},
			insert: []time.Duration{time.Second},
			exp:    []time.Duration{time.Second},
		},
		{
			name:   "empty input, multiple inserts",
			arr:    []time.Duration{},
			insert: []time.Duration{time.Second, 3 * time.Second, 2 * time.Second},
			exp:    []time.Duration{time.Second, 2 * time.Second, 3 * time.Second},
		},
		{
			name:   "non-empty input, multiple inserts",
			arr:    []time.Duration{time.Second, 3 * time.Second, 4 * time.Second, 9 * time.Second},
			insert: []time.Duration{time.Second, 7 * time.Second, 2 * time.Second},
			exp:    []time.Duration{time.Second, time.Second, 2 * time.Second, 3 * time.Second, 4 * time.Second, 7 * time.Second, 9 * time.Second},
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			for _, d := range v.insert {
				v.arr = sortedAppend(v.arr, d)
			}

			testutil.Equals(t, v.exp, v.arr)
		})
	}
}
