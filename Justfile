set shell := ["bash", "-euo", "pipefail", "-c"]

build:
    go build -o relinted ./cmd/relinted/

test:
    go test ./internal/... ./cmd/...

run:
    go run ./cmd/relinted/

lint:
    go vet ./internal/... ./cmd/relinted/

format:
    gofmt -w .
