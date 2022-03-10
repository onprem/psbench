SHELL=/usr/bin/env bash -o pipefail

BIN_NAME ?= psbench

default: $(BIN_NAME)
all: clean $(BIN_NAME)

$(BIN_NAME): main.go $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://proxy.golang.org go build -a -ldflags '-s -w' -o $@ .

.PHONY: build
build: $(BIN_NAME)

.PHONY: clean
clean:
	-rm $(BIN_NAME)
