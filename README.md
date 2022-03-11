# `psbench`

A tool to benchmark PromQL query performance across multiple workers/clients against a Promscale instance.

```bash mdox-exec="./psbench -help"
Usage of ./psbench:
  -log.format string
    	The log format to use. Options: 'logfmt', 'json'. (default "logfmt")
  -log.level string
    	The log filtering level. Options: 'error', 'warn', 'info', 'debug'. (default "info")
  -promscale.url string
    	The URL of the Promscale instance to run the benchmark against.
  -queries.file string
    	Path to CSV file that contains rows with queries to execute. Format: PromQL query,start_time,end_time,step size.
  -workers int
    	The number of workers/clients to run parallelly to query the Promscale instance. (default 1)
```
