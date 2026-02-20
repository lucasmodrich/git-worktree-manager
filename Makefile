.PHONY: build test install clean fmt lint

# Binary name (produces gwtm.exe on Windows automatically)
BINARY=gwtm

# Build info injected at compile time.
# VERSION uses the exact release tag if HEAD is tagged, otherwise "dev".
# COMMIT and DATE always identify the exact build.
VERSION := $(shell git describe --exact-match --tags 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")

LDFLAGS := -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

# Build the binary
build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) ./cmd/git-worktree-manager

# Run all tests
test:
	go test ./...

# Install binary to GOPATH/bin or HOME/.git-worktree-manager
install: build
	@if [ -n "$$GOPATH" ]; then \
		cp $(BINARY) $$GOPATH/bin/; \
	else \
		mkdir -p $$HOME/.git-worktree-manager; \
		cp $(BINARY) $$HOME/.git-worktree-manager/; \
	fi

# Clean build artifacts
clean:
	rm -f $(BINARY) $(BINARY).exe
	go clean

# Format code
fmt:
	gofmt -s -w .

# Run linter (if golangci-lint is installed)
lint:
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"
