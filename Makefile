BINARY := iluvatar
DIST   := .dist

MISC := misc

# Version of Go to install when the machine doesn't have one.
GO_VERSION      ?= 1.27.1
GO_INSTALL_DIR  ?= /usr/local/go
GOPROXY         ?= https://proxy.golang.org,direct

export GO_VERSION
export GO_INSTALL_DIR
export GOPROXY

.PHONY: build install test fmt vet lint vulns clean setup setup-go setup-tools

build:
	go build -o $(DIST)/$(BINARY) ./cmd/iluvatar

install:
	go install ./cmd/iluvatar

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

vulns:
	govulncheck ./...

clean:
	rm -rf $(DIST)

# setup installs every technology needed to work on this project: the Go
# toolchain (if missing), plus the linter and vulnerability scanner used by
# `make lint` and `make vulns`. Run it once after cloning in a new
# environment.
setup: setup-tools

setup-go:
	@$(MISC)/setup-go.sh

setup-tools: setup-go
	@$(MISC)/setup-tools.sh
