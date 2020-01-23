SHELL := bash
.ONESHELL:

GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

get:
	@$(GOGET) github.com/sparrc/go-ping
	@$(GOGET) github.com/influxdata/influxdb1-client/v2

format:
	$(GOFMT) ./...

build/amd64:
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/linux-amd64/pingserv cmd/pingserv/main.go

build: build/amd64

clean:
	@rm -rf bin

all: dir get format build
