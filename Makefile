SHELL := bash
.ONESHELL:

GO=$(shell which go)
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

get:
	@$(GOGET) github.com/stenya/go-ping
	@$(GOGET) github.com/influxdata/influxdb1-client/v2
	@$(GOGET) github.com/spf13/viper

format:
	$(GOFMT) ./...

build/linux-amd64:
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/linux-amd64/pingflux cmd/pingflux/main.go

build/darwin-amd64:
	export GOOS=darwin
	export GOARCH=amd64
	$(GOBUILD) -o bin/darwin-amd64/pingflux cmd/pingflux/main.go

build: build/linux-amd64 build/darwin-amd64

clean:
	@rm -rf bin

all: dir get format build
