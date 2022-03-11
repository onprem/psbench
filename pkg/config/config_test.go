package config

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/efficientgo/tools/core/pkg/testutil"
)

func TestParseQueriesFile(t *testing.T) {
	cases := []struct {
		name   string
		input  io.Reader
		exp    []Query
		expErr error
	}{
		{
			name:   "empty input",
			input:  strings.NewReader(""),
			exp:    []Query{},
			expErr: nil,
		},
		{
			name: "given sample",
			input: strings.NewReader(`demo_cpu_usage_seconds_total{mode="idle"}|1597056698698|1597059548699|15000
avg by(instance) (demo_cpu_usage_seconds_total)|1597057698698|1597058548699|60000
avg without(instance, mode) (demo_cpu_usage_seconds_total)|1597056698698|1597059548699|120000
`),
			exp: []Query{
				{`demo_cpu_usage_seconds_total{mode="idle"}`, time.Unix(1597056698698, 0), time.Unix(1597059548699, 0), 15000 * time.Second},
				{`avg by(instance) (demo_cpu_usage_seconds_total)`, time.Unix(1597057698698, 0), time.Unix(1597058548699, 0), 60000 * time.Second},
				{`avg without(instance, mode) (demo_cpu_usage_seconds_total)`, time.Unix(1597056698698, 0), time.Unix(1597059548699, 0), 120000 * time.Second},
			},
			expErr: nil,
		},
		{
			name: "invalid input",
			input: strings.NewReader(`demo_cpu_usage_seconds_total{mode="idle"}|1597056698698|1597059548699
avg by(instance) (demo_cpu_usage_seconds_total)|1597057698698|1597058548699`),
			exp:    []Query{},
			expErr: errInvalidRow,
		},
	}

	for _, v := range cases {
		t.Run(v.name, func(t *testing.T) {
			got, err := parseQueriesFile(v.input)

			testutil.Equals(t, v.expErr, err)

			testutil.Equals(t, v.exp, got)
		})
	}
}
