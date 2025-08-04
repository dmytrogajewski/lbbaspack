# LBaaS Packet Catcher - ECS Edition Makefile
# Based on Go project layout guidelines

# Project information
PROJECT_NAME := lbbaspack
BINARY_NAME := lbbaspack
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# Go parameters
GOCMD := go
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOVET := $(GOCMD) vet
GOFMT := $(GOCMD) fmt
GOLINT := golangci-lint

# Build flags
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

# Directories
BIN_DIR := bin
DIST_DIR := dist
BUILD_DIR := build
COVERAGE_DIR := coverage

# Files
MAIN_FILE := main.go
BINARY := $(BIN_DIR)/$(BINARY_NAME)
DIST_BINARY := $(DIST_DIR)/$(BINARY_NAME)-$(VERSION)

# Test parameters
TEST_FLAGS := -v -race -cover
COVERAGE_FILE := $(COVERAGE_DIR)/coverage.out
COVERAGE_HTML := $(COVERAGE_DIR)/coverage.html

# Development parameters
WATCH_DIRS := . engine/ ui/ config/ utils/
WATCH_EXTENSIONS := go

# Default target
.DEFAULT_GOAL := help

# Create necessary directories
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(DIST_DIR):
	mkdir -p $(DIST_DIR)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(COVERAGE_DIR):
	mkdir -p $(COVERAGE_DIR)

# Help target
.PHONY: help
help: ## Show this help message
	@echo "LBaaS Packet Catcher - ECS Edition"
	@echo "=================================="
	@echo ""
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
	@echo ""
	@echo "Examples:"
	@echo "  make build     # Build the application"
	@echo "  make test      # Run all tests"
	@echo "  make run       # Run the application"
	@echo "  make clean     # Clean build artifacts"

# Build targets
.PHONY: build
build: $(BIN_DIR) ## Build the application
	@echo "Building $(BINARY_NAME) version $(VERSION)..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY) $(MAIN_FILE)
	@echo "Build complete: $(BINARY)"

.PHONY: build-release
build-release: $(DIST_DIR) ## Build release version with optimizations
	@echo "Building release version $(VERSION)..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(DIST_BINARY)-linux-amd64 $(MAIN_FILE)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(DIST_BINARY)-darwin-amd64 $(MAIN_FILE)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -a -installsuffix cgo -o $(DIST_BINARY)-windows-amd64.exe $(MAIN_FILE)
	@echo "Release builds complete in $(DIST_DIR)"

.PHONY: build-debug
build-debug: $(BIN_DIR) ## Build debug version with symbols
	@echo "Building debug version..."
	$(GOBUILD) $(LDFLAGS) -gcflags="all=-N -l" -o $(BINARY)-debug $(MAIN_FILE)
	@echo "Debug build complete: $(BINARY)-debug"

# Test targets
.PHONY: test
test: $(COVERAGE_DIR) ## Run all tests
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) ./...

.PHONY: test-short
test-short: ## Run tests without race detection
	@echo "Running short tests..."
	$(GOTEST) -v -cover ./...

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo "Running verbose tests..."
	$(GOTEST) -v -race -cover -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

.PHONY: test-benchmark
test-benchmark: ## Run benchmark tests
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

.PHONY: test-coverage
test-coverage: $(COVERAGE_DIR) ## Generate test coverage report
	@echo "Generating coverage report..."
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

# Run targets
.PHONY: run
run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY)

.PHONY: run-debug
run-debug: build-debug ## Build and run debug version
	@echo "Running debug version..."
	./$(BINARY)-debug

.PHONY: run-dev
run-dev: ## Run in development mode (with hot reload if available)
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_FILE)

# Development targets
.PHONY: dev
dev: ## Start development environment
	@echo "Starting development environment..."
	@echo "Available commands:"
	@echo "  make run-dev     # Run with go run"
	@echo "  make test        # Run tests"
	@echo "  make lint        # Run linter"
	@echo "  make fmt         # Format code"
	@echo "  make deps        # Install dependencies"

.PHONY: watch
watch: ## Watch for changes and run tests (requires fswatch)
	@echo "Watching for changes..."
	@if command -v fswatch >/dev/null 2>&1; then \
		fswatch -o $(WATCH_DIRS) | xargs -n1 -I{} make test-short; \
	else \
		echo "fswatch not found. Install it or use 'make test' manually."; \
	fi

# Code quality targets
.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting code..."
	gofmt -s -w .

.PHONY: fmt-check
fmt-check: ## Check if code is properly formatted
	@echo "Checking code formatting..."
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Code is not formatted. Run 'make fmt' to fix."; \
		exit 1; \
	else \
		echo "Code is properly formatted."; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

.PHONY: lint
lint: ## Run golangci-lint
	@echo "Running linter..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo "Running linter with auto-fix..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run --fix; \
	else \
		echo "golangci-lint not found. Install it with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Dependency management
.PHONY: deps
deps: ## Install dependencies
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

.PHONY: deps-update
deps-update: ## Update dependencies
	@echo "Updating dependencies..."
	$(GOMOD) get -u ./...
	$(GOMOD) tidy

.PHONY: deps-check
deps-check: ## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	$(GOMOD) list -u

# Clean targets
.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	$(GOCLEAN)
	rm -rf $(BIN_DIR)
	rm -rf $(BUILD_DIR)
	@echo "Clean complete"

.PHONY: clean-all
clean-all: clean ## Clean all artifacts including coverage and dist
	@echo "Cleaning all artifacts..."
	rm -rf $(COVERAGE_DIR)
	rm -rf $(DIST_DIR)
	@echo "All artifacts cleaned"

# Documentation targets
.PHONY: docs
docs: ## Generate documentation
	@echo "Generating documentation..."
	@if command -v godoc >/dev/null 2>&1; then \
		echo "Starting godoc server at http://localhost:6060"; \
		godoc -http=:6060; \
	else \
		echo "godoc not found. Install it with: go install golang.org/x/tools/cmd/godoc@latest"; \
	fi

# Security targets
.PHONY: security
security: ## Run security checks
	@echo "Running security checks..."
	@if command -v gosec >/dev/null 2>&1; then \
		gosec ./...; \
	else \
		echo "gosec not found. Install it with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Performance targets
.PHONY: profile
profile: build ## Generate performance profile
	@echo "Generating performance profile..."
	./$(BINARY) -cpuprofile=cpu.prof -memprofile=mem.prof
	@echo "Profiles generated: cpu.prof, mem.prof"

.PHONY: profile-cpu
profile-cpu: ## Analyze CPU profile
	@if [ -f cpu.prof ]; then \
		$(GOCMD) tool pprof cpu.prof; \
	else \
		echo "CPU profile not found. Run 'make profile' first."; \
	fi

.PHONY: profile-mem
profile-mem: ## Analyze memory profile
	@if [ -f mem.prof ]; then \
		$(GOCMD) tool pprof mem.prof; \
	else \
		echo "Memory profile not found. Run 'make profile' first."; \
	fi

# CI/CD targets
.PHONY: ci
ci: fmt-check vet lint test ## Run CI checks
	@echo "CI checks completed successfully"

.PHONY: ci-full
ci-full: ci test-coverage security ## Run full CI pipeline
	@echo "Full CI pipeline completed successfully"

# Release targets
.PHONY: release
release: clean-all test-coverage build-release ## Create release build
	@echo "Release $(VERSION) created successfully"
	@echo "Binaries available in $(DIST_DIR)"

.PHONY: release-check
release-check: ## Check release readiness
	@echo "Checking release readiness..."
	@make ci-full
	@make build-release
	@echo "Release check completed"

# Docker targets (if applicable)
.PHONY: docker-build
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(PROJECT_NAME):$(VERSION) .
	docker tag $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest

.PHONY: docker-run
docker-run: ## Run Docker container
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(PROJECT_NAME):$(VERSION)

.PHONY: docker-clean
docker-clean: ## Clean Docker images
	@echo "Cleaning Docker images..."
	docker rmi $(PROJECT_NAME):$(VERSION) $(PROJECT_NAME):latest 2>/dev/null || true

# Utility targets
.PHONY: version
version: ## Show version information
	@echo "Project: $(PROJECT_NAME)"
	@echo "Version: $(VERSION)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Git Commit: $(GIT_COMMIT)"

.PHONY: info
info: ## Show project information
	@echo "Project Information:"
	@echo "  Name: $(PROJECT_NAME)"
	@echo "  Binary: $(BINARY_NAME)"
	@echo "  Version: $(VERSION)"
	@echo "  Go Version: $(shell $(GOCMD) version)"
	@echo "  Architecture: $(shell $(GOCMD) env GOOS)/$(shell $(GOCMD) env GOARCH)"

.PHONY: check
check: ## Check if all tools are available
	@echo "Checking required tools..."
	@command -v $(GOCMD) >/dev/null 2>&1 || { echo "Go is not installed"; exit 1; }
	@echo "✓ Go is installed"
	@command -v git >/dev/null 2>&1 || { echo "Git is not installed"; exit 1; }
	@echo "✓ Git is installed"
	@echo "All required tools are available"

# Install targets
.PHONY: install
install: build ## Install the application
	@echo "Installing $(BINARY_NAME)..."
	cp $(BINARY) /usr/local/bin/$(BINARY_NAME)
	@echo "Installation complete"

.PHONY: uninstall
uninstall: ## Uninstall the application
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "Uninstallation complete"

# Backup and restore targets
.PHONY: backup
backup: ## Create backup of current state
	@echo "Creating backup..."
	@mkdir -p backups
	@tar -czf backups/$(PROJECT_NAME)-$(VERSION)-$(shell date +%Y%m%d-%H%M%S).tar.gz \
		--exclude='.git' \
		--exclude='bin' \
		--exclude='build' \
		--exclude='dist' \
		--exclude='coverage' \
		--exclude='*.prof' \
		--exclude='backups' \
		.
	@echo "Backup created"

# Default development workflow
.PHONY: dev-setup
dev-setup: deps check ## Setup development environment
	@echo "Development environment setup complete"
	@echo "Run 'make dev' to start development"

.PHONY: dev-workflow
dev-workflow: fmt lint test ## Run development workflow
	@echo "Development workflow completed" 