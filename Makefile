SHELL := bash
.ONESHELL:

VER=$(shell git describe --tags)
GO=$(shell which go)
GOMOD=$(GO) mod
GOFMT=$(GO) fmt
GOBUILD=$(GO) build -mod=readonly -ldflags "-X main.version=$(VER)"

dir:
	@if [ ! -d bin ]; then mkdir -p bin; fi

mod:
	@$(GOMOD) download

format:
	@$(GOFMT) ./...

build/linux/armv7: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=arm
	export GOARM=7
	$(GOBUILD) -o bin/pingflux-linux-$(VER:v%=%)-armv7 cmd/pingflux/main.go

build/linux/arm64: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=arm64
	$(GOBUILD) -o bin/pingflux-linux-$(VER:v%=%)-arm64 cmd/pingflux/main.go

build/linux/386: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=386
	$(GOBUILD) -o bin/pingflux-linux-$(VER:v%=%)-386 cmd/pingflux/main.go

build/linux/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=linux
	export GOARCH=amd64
	$(GOBUILD) -o bin/pingflux-linux-$(VER:v%=%)-amd64 cmd/pingflux/main.go

build/linux: build/linux/armv7 build/linux/arm64 build/linux/386 build/linux/amd64

build/darwin/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=darwin
	export GOARCH=amd64
	$(GOBUILD) -o bin/pingflux-darwin-$(VER:v%=%)-amd64 cmd/pingflux/main.go

build/darwin: build/darwin/amd64

build/windows/amd64: dir mod
	export CGO_ENABLED=0
	export GOOS=windows
	export GOARCH=amd64
	$(GOBUILD) -o bin/pingflux-windows-$(VER:v%=%)-amd64 cmd/pingflux/main.go

build/windows: build/windows/amd64

build: build/linux build/darwin build/windows

clean:
	@rm -rf bin

all: format build
