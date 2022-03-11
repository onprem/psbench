include .bingo/Variables.mk

SHELL=/usr/bin/env bash -o pipefail

BIN_NAME ?= psbench
MDOX_VALIDATE_CONFIG ?= .mdox.validate.yaml

default: $(BIN_NAME)
all: clean $(BIN_NAME)

$(BIN_NAME): main.go $(wildcard *.go) $(wildcard */*.go)
	CGO_ENABLED=0 GO111MODULE=on GOPROXY=https://proxy.golang.org go build -a -ldflags '-s -w' -o $@ .

.PHONY: build
build: $(BIN_NAME)

.PHONY: docs
docs: build $(MDOX) ## Generates config snippets and doc formatting.
	@echo ">> generating docs $(PATH)"
	PATH=${PATH}:$(GOBIN) $(MDOX) fmt -l --links.validate.config-file=$(MDOX_VALIDATE_CONFIG) *.md

.PHONY: clean
clean:
	-rm $(BIN_NAME)
