# Makefile for the Automater CLI

# Variables
BINARY_NAME=automater
GO=go
GOFMT=gofmt
GO_FILES=$(shell find . -name '*.go' | grep -v /vendor/)
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X github.com/moabdelazem/automater/cmd/root.Version=$(VERSION)"

# Default target
.PHONY: all
all: fmt test build

# Build the binary
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@$(GO) build -o $(BINARY_NAME) $(LDFLAGS)

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@$(GO) test ./...

# Format all code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@$(GOFMT) -w $(GO_FILES)

# Install the binary
.PHONY: install
install:
	@echo "Installing $(BINARY_NAME)..."
	@$(GO) install $(LDFLAGS)

# Clean up 
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)
	@rm -rf dist/

# Create a dist directory with binaries for multiple platforms
.PHONY: dist
dist: clean
	@echo "Creating distribution packages..."
	@mkdir -p dist
	
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)_linux_amd64
	
	@echo "Building for MacOS..."
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)_darwin_amd64
	
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o dist/$(BINARY_NAME)_windows_amd64.exe

# Run the binary
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BINARY_NAME)

# Help
.PHONY: help
help:
	@echo "Automater CLI Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make              Build the application after formatting code and running tests"
	@echo "  make build        Build the binary"
	@echo "  make test         Run tests"
	@echo "  make fmt          Format the code"
	@echo "  make install      Install the binary"
	@echo "  make clean        Clean up binary and dist directory"
	@echo "  make dist         Create distribution packages for multiple platforms"
	@echo "  make run          Build and run the binary"
	@echo "  make help         Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION           Set the version for the build (default: git tag or 'dev')"