.PHONY: build test lint clean install run docker-build docker-run

# Build variables
BINARY_NAME=gomailserver
BUILD_DIR=./build
MAIN_PATH=./cmd/gomailserver
VERSION?=dev
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -s -w"

# Go commands
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOGET=$(GO) get
GOMOD=$(GO) mod
GOLINT=golangci-lint

all: clean lint test build

build: build-ui
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

build-ui:
	@echo "Building admin UI..."
	@cd web/admin && npm install && npm run build
	@echo "Admin UI build complete"

build-static: build-ui
	@echo "Building static binary..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 $(GOBUILD) $(LDFLAGS) -tags 'osusergo netgo static_build' \
		-a -installsuffix cgo -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Static build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test:
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	@echo "Tests complete"

test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

lint:
	@echo "Running linter..."
	$(GOLINT) run
	@echo "Lint complete"

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

install: build
	@echo "Installing $(BINARY_NAME)..."
	@install -d /usr/local/bin
	@install -m 755 $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@install -d /etc/gomailserver
	@echo "Installation complete"

run: build
	@echo "Running $(BINARY_NAME)..."
	$(BUILD_DIR)/$(BINARY_NAME)

docker-build:
	@echo "Building Docker image..."
	docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .
	@echo "Docker image built"

docker-run:
	@echo "Running Docker container..."
	docker run -p 25:25 -p 143:143 -p 465:465 -p 587:587 -p 993:993 \
		-v gomailserver-data:/data \
		$(BINARY_NAME):latest

.DEFAULT_GOAL := build
