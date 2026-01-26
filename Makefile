# Go build settings
GO ?= go
GOBUILD ?= $(GO) build
GOTEST ?= $(GO) test
GOFMT ?= $(GO) fmt
GOVET ?= $(GO) vet
BUILD_DIR ?= build
BINARY_NAME ?= gomailserver
DOCKER_IMAGE ?= gomailserver:latest
DOCKER_TAG ?= gomailserver:test

.PHONY: build build-static test test-coverage lint clean docker-build docker-run deps bench-perf release-snapshot release-check

# Default target
all: build test

# Build the main application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/gomailserver
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

# Build static binary (for Docker)
build-static:
	@echo "Building static $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) \
		-ldflags='-w -s -extldflags "-static"' \
		-a \
		-installsuffix cgo \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		./cmd/gomailserver
	@echo "Built static $(BUILD_DIR)/$(BINARY_NAME)"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests completed. Coverage: coverage.out"

# Run tests with coverage
test-coverage: test
	@echo "Generating coverage report..."
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	@echo "Clean completed"

# Download/update dependencies
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	@echo "Dependencies updated"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-build-test:
	@echo "Building Docker test image..."
	docker build -t $(DOCKER_TAG) .
	@echo "Docker test image built: $(DOCKER_TAG)"

docker-run:
	@echo "Starting gomailserver with Docker Compose..."
	docker-compose up -d

docker-stop:
	@echo "Stopping gomailserver..."
	docker-compose down

docker-logs:
	@echo "Following Docker logs..."
	docker-compose logs -f

# Performance benchmark
bench-perf:
	@echo "Building performance benchmark tool..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/bench-perf ./cmd/bench-perf
	@echo "Performance benchmark tool built: $(BUILD_DIR)/bench-perf"

# Install binary system-wide
install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	sudo chmod 755 /usr/local/bin/$(BINARY_NAME)
	@echo "Installed to /usr/local/bin/$(BINARY_NAME)"

# Development server
dev: build
	@echo "Starting development server..."
	./$(BUILD_DIR)/$(BINARY_NAME) run --config gomailserver.yaml

# GoReleaser targets
release-snapshot:
	@echo "Building snapshot release (no publish)..."
	goreleaser release --snapshot --clean
	@echo "Snapshot built in ./dist/"

release-check:
	@echo "Checking GoReleaser configuration..."
	goreleaser check

# Help target
help:
	@echo "Available targets:"
	@echo "  build         - Build the main application"
	@echo "  build-static  - Build static binary for Docker"
	@echo "  test          - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  lint          - Run linter"
	@echo "  clean         - Clean build artifacts"
	@echo "  deps          - Download/update dependencies"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Start with Docker Compose"
	@echo "  docker-stop    - Stop Docker Compose services"
	@echo "  docker-logs   - Follow Docker logs"
	@echo "  bench-perf    - Build performance benchmark"
	@echo "  install       - Install binary system-wide"
	@echo "  dev           - Start development server"
	@echo "  release-snapshot - Build snapshot release (no publish)"
	@echo "  release-check    - Check GoReleaser configuration"
	@echo "  help          - Show this help"