SHELL := bash
.ONESHELL:

GO=$(shell which go)
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOBUILD=$(GO) build

dir:
	@if [ ! -d bin ] ; then mkdir -p bin ; fi

mod:
	@$(GOMOD) download

format:
	$(GOFMT) ./...

build/linux/amd64: dir mod
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/linux-amd64/pingflux cmd/pingflux/main.go

build/linux/arm: dir mod
	export GOOS=linux
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/linux-arm/pingflux cmd/pingflux/main.go

build/linux/arm64: dir mod
	export GOOS=linux
	export GOARCH=arm64
	$(GOBUILD) -o bin/linux-arm64/pingflux cmd/pingflux/main.go

build/linux: build/linux/amd64 build/linux/arm build/linux/arm64

build/darwin/amd64: dir mod
	export GOOS=darwin
	export GOARCH=amd64
	$(GOBUILD) -o bin/darwin-amd64/pingflux cmd/pingflux/main.go

build/darwin: build/darwin/amd64

build/windows/amd64: dir mod
	export GOOS=windows
	export GOARCH=amd64
	$(GOBUILD) -o bin/windows-amd64/pingflux cmd/pingflux/main.go

build/windows: build/windows/amd64

build: build/linux build/darwin build/windows

clean:
	@rm -rf bin

all: format build
