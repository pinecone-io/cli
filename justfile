# Load environment variables from .env before recipes run
set dotenv-load := true

default:
    @just --list

# Check to see if Go is available locally and error with command-not-found if not
ensure-go:
    @if bin="$(command -v go)"; then \
        echo "Found Go at: $bin"; \
        "$bin" version; \
    else \
        echo "Go not found in PATH"; \
        exit 127; \
    fi

# Check to see if goreleaser is available locally and error with command-not-found if not
ensure-goreleaser:
    @if bin="$(command -v goreleaser)"; then \
        echo "Found goreleaser at: $bin"; \
        "$bin" --version; \
    else \
        echo "goreleaser not found in PATH"; \
        exit 127; \
    fi

# Run all tests for the CLI
test *ARGS: ensure-go
    go test -v ./... {{ARGS}}

# Generate man pages for the CLI, output in ./man
gen-manpages *ARGS: ensure-go
    go run cmd/gen-manpages/main.go {{ARGS}}

# Build the CLI binary locally using goreleaser: current OS, artifacts in ./dist
# 
build: ensure-go ensure-goreleaser
    goreleaser build --single-target --snapshot --clean

# Build the CLI binary locally using goreleaser: all supported OSes (defined in.goreleaser.yaml), artifacts in ./dist
build-all: ensure-go ensure-goreleaser
    goreleaser build --snapshot --clean