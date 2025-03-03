.PHONY: all build clean test test-coverage test-integration test-all

# Binary name
BINARY_NAME=gogo
# Binary directory
BIN_DIR=./bin

# Go commands
GO ?= go
GOBUILD = $(GO) build
GOCLEAN = $(GO) clean
GOTEST = $(GO) test
GOGET = $(GO) get

# Version info from git
GIT_COMMIT=$(shell git rev-parse --short HEAD || echo "unknown")
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+DIRTY" || echo "")
GIT_TAG=$(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')

# Get the module name from go.mod
MODULE_NAME=$(shell grep "^module" go.mod | awk '{print $$2}')

# Linker flags
LDFLAGS=-ldflags "-X $(MODULE_NAME)/cmd/gogo.Version=$(GIT_TAG) \
-X $(MODULE_NAME)/cmd/gogo.Commit=$(GIT_COMMIT)$(GIT_DIRTY) \
-X $(MODULE_NAME)/cmd/gogo.BuildDate=$(BUILD_DATE)"

# Default target (run tests and build binary)
all:
	@echo "Running all tests and building..."
	$(MAKE) test-all
	$(MAKE) build
	@echo "All done!"

# Build binary
build:
	@echo "Building $(BINARY_NAME)..."
	@echo "Git commit: $(GIT_COMMIT)$(GIT_DIRTY)"
	@echo "Git tag: $(GIT_TAG)"
	@echo "Build date: $(BUILD_DATE)"
	@mkdir -p $(BIN_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME)
	@echo "Build complete: $(BIN_DIR)/$(BINARY_NAME)"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@$(GOCLEAN)
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...
	@echo "Tests complete"

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v ./... -coverprofile=coverage.out
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated at coverage.html"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	GOGO_INTEGRATION_TEST=1 $(GOTEST) -v ./test/integration/
	@echo "Integration tests complete"

# Run all tests (unit and integration) but continue even if tests fail
test-all:
	@echo "Running all tests..."
	-$(GOTEST) -v ./...
	-GOGO_INTEGRATION_TEST=1 $(GOTEST) -v ./test/integration/
	@echo "All tests complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...
	@echo "Dependencies installed"

# Lint the code
lint:
	@echo "Linting code..."
	golangci-lint run ./...
	@echo "Lint complete"

# Format the code
fmt:
	@echo "Formatting code..."
	gofmt -s -w .
	@echo "Format complete"

# Install pre-commit hooks
install-hooks:
	@echo "Installing pre-commit hooks..."
	pre-commit install
	@echo "Pre-commit hooks installed"

# Test Homebrew formula locally
test-brew: build
	@echo "Testing Homebrew formula locally..."
	@echo "1. Building with goreleaser in snapshot mode..."
	goreleaser release --snapshot --rm-dist --skip-publish
	@echo "2. Creating local tap..."
	rm -rf /tmp/homebrew-gogo || true
	mkdir -p /tmp/homebrew-gogo
	cp homebrew-gogo/gogo.rb /tmp/homebrew-gogo/
	@echo "3. Updating formula with local paths..."
	dist_path=$$(pwd)/dist
	darwin_amd64_path=$$(find $$dist_path -name "gogo_*_darwin_amd64.tar.gz" | head -n 1)
	darwin_amd64_sha=$$(shasum -a 256 $$darwin_amd64_path | cut -d ' ' -f 1)
	sed -i '' "s|url \".*darwin_amd64.tar.gz\"|url \"file://$$darwin_amd64_path\"|" /tmp/homebrew-gogo/gogo.rb
	sed -i '' "s|sha256 \".*\" # Replace with.*darwin_amd64|sha256 \"$$darwin_amd64_sha\" # darwin_amd64|" /tmp/homebrew-gogo/gogo.rb
	@echo "4. Installing from local tap..."
	brew tap --force gogo-local /tmp/homebrew-gogo
	brew uninstall --force gogo || true
	brew install --verbose --debug gogo-local/gogo
	@echo "5. Testing installation..."
	gogo version
	@echo "Homebrew formula test complete!"

# Help target
help:
	@echo "Available targets:"
	@echo "  all               - Run all tests and build the binary"
	@echo "  build             - Build the binary to $(BIN_DIR)/$(BINARY_NAME)"
	@echo "  clean             - Clean build artifacts"
	@echo "  test              - Run tests"
	@echo "  test-coverage     - Run tests with coverage reporting"
	@echo "  test-integration  - Run integration tests"
	@echo "  test-all          - Run both unit and integration tests"
	@echo "  deps              - Install dependencies"
	@echo "  lint              - Lint the code"
	@echo "  fmt               - Format the code"
	@echo "  install-hooks     - Install pre-commit hooks"
	@echo "  test-brew         - Test Homebrew formula locally"
