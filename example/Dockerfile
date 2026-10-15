FROM golang:1.27.1-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/example ./cmd/example

FROM alpine:3.20
RUN adduser -D -u 10001 app
COPY --from=builder /out/example /usr/local/bin/example
USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/example"]
