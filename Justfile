set shell := ["bash", "-euo", "pipefail", "-c"]

build:
    go build -ldflags "-X main.version=$(git describe --tags)" -o relinted ./cmd/relinted/

test:
    go test ./internal/... ./cmd/...

test-c:
    go test ./internal/formatter/... ./internal/tokenizer/... ./internal/io/...

test-perl:
    go test ./internal/perl/...

test-rust:
    go test ./internal/rust/...

test-js:
    go test ./internal/js/...

test-go:
    go test ./internal/go/...

test-java:
    go test ./internal/java/...

test-php:
    go test ./internal/php/...

test-swift:
    go test ./internal/swift/...

test-ts:
    go test ./internal/js/...

test-cs:
    go test ./internal/java/...

run:
    go run ./cmd/relinted/

lint:
    go vet ./internal/... ./cmd/relinted/

format:
    gofmt -w .
