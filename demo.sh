#!/bin/bash

docker-compose up -d

sleep 5

curl -v \
  -H "Content-Type: application/x-protobuf" \
  -H "Content-Encoding: snappy" \
  -H "X-Prometheus-Remote-Write-Version: 0.1.0" \
  --data-binary "@docker/real-dataset.sz" \
  http://localhost:9201/write

sleep 2

make build

echo -e "\n--------------------------------------------\n"

./psbench -workers 3 -queries.file ./docker/obs-queries.csv -promscale.url http://localhost:9201
