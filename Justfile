set shell := ["bash", "-euo", "pipefail", "-c"]

build:
    go build -o relinted ./cmd/relinted/

test:
    go test ./internal/... ./cmd/...

test-c:
    go test ./internal/formatter/... ./internal/tokenizer/... ./internal/io/...

test-perl:
    go test ./internal/perl/...

run:
    go run ./cmd/relinted/

lint:
    go vet ./internal/... ./cmd/relinted/

format:
    gofmt -w .
