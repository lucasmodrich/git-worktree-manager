.PHONY: build test install clean fmt lint

# Binary name
BINARY=gwtm

# Build the binary
build:
	go build -o $(BINARY) ./cmd/git-worktree-manager

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
	rm -f $(BINARY)
	go clean

# Format code
fmt:
	gofmt -s -w .

# Run linter (if golangci-lint is installed)
lint:
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run || echo "golangci-lint not installed, skipping"
