BINARY := manwe
DIST   := .dist

.PHONY: build run test fmt vet lint vulns proto migrate-up migrate-down docker-build clean

build:
	go build -o $(DIST)/$(BINARY) ./cmd/$(BINARY)

run:
	go run ./cmd/$(BINARY)

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

# proto regenerates the *.pb.go/*_grpc.pb.go files from every resource's
# .proto. Requires protoc, protoc-gen-go, and protoc-gen-go-grpc on PATH;
# `iluvatar generate` installs whichever of them are missing.
proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $$(find internal -name '*.proto')

# migrate-up/migrate-down apply/revert migrations/ with golang-migrate's
# `migrate` CLI, which isn't installed for you. DATABASE_DSN must be set.
migrate-up:
	migrate -path migrations -database "$$DATABASE_DSN" up

migrate-down:
	migrate -path migrations -database "$$DATABASE_DSN" down

docker-build:
	docker build -t $(BINARY) .

clean:
	rm -rf $(DIST)
