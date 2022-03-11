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
    	Path to CSV file that contains rows with queries to execute. Format: PromQL query|start_time|end_time|step size.
  -workers int
    	The number of workers/clients to run parallelly to query the Promscale instance. (default 1)
```

## Usage

> NOTE: Instead of running all these steps manually, you can run the `demo.sh` script.

- You'll need `docker` running and `docker-compose` to run TimescaleDB and Promscale. Run `docker-compose up -d`, to start everything using the provided compose file.
- Now we need to write the example dataset. Run

  ```bash
  curl -v \
    -H "Content-Type: application/x-protobuf" \
    -H "Content-Encoding: snappy" \
    -H "X-Prometheus-Remote-Write-Version: 0.1.0" \
    --data-binary "@docker/real-dataset.sz" \
    http://localhost:9201/write
  ```
- Build the `psbench` binary by running `make build`.
- Run the benchmark

  ```bash
  ./psbench -workers 3 -queries.file ./docker/obs-queries.csv -promscale.url http://localhost:9201
  ```
- After you are done, don't forget to tear down the docker compose stack by running `docker-compose down`.

## Example output

```
$ ./psbench -workers 3 -queries.file ./docker/obs-queries.csv -promscale.url http://localhost:9201

Total Number of Queries: 	10
Total Processing Time: 		24.681721ms
Minimum Query Time: 		1.131087ms
Maximum Query Time: 		3.814267ms
Median Query Time: 			2.277084ms
Avergae Query Time: 		2.468172ms
```
