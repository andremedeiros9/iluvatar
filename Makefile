BINARY := iluvatar
DIST   := .dist

.PHONY: build install test fmt vet lint vulns clean

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
