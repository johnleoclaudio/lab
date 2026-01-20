.PHONY: help build run test clean docker-up docker-down migrate-up migrate-down sqlc-gen lint fmt

# Variables
APP_NAME=go-starter
MAIN_PATH=./cmd/server
MIGRATIONS_PATH=./migrations
DATABASE_URL?=postgres://postgres:postgres@localhost:5432/go_starter?sslmode=disable

# Default target
help:
	@echo "Available commands:"
	@echo "  make build        - Build the application"
	@echo "  make run          - Run the application"
	@echo "  make test         - Run tests with coverage"
	@echo "  make test-unit    - Run unit tests only"
	@echo "  make test-integration - Run integration tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-up    - Start Docker containers"
	@echo "  make docker-down  - Stop Docker containers"
	@echo "  make migrate-up   - Run database migrations"
	@echo "  make migrate-down - Rollback database migrations"
	@echo "  make sqlc-gen     - Generate sqlc code"
	@echo "  make lint         - Run linter"
	@echo "  make fmt          - Format code"
	@echo "  make deps         - Download dependencies"

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o bin/$(APP_NAME) $(MAIN_PATH)
	@echo "Build complete: bin/$(APP_NAME)"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run $(MAIN_PATH)

# Run tests with coverage
test:
	@echo "Running tests with coverage..."
	@go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
	@go tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run unit tests only
test-unit:
	@echo "Running unit tests..."
	@go test -v -race -short ./...

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -race -run Integration ./tests/integration/...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.txt coverage.html
	@echo "Clean complete"

# Start Docker containers
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Containers started"

# Stop Docker containers
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Containers stopped"

# Run database migrations up
migrate-up:
	@echo "Running migrations..."
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up; \
	else \
		echo "Error: migrate tool not installed. Install with:"; \
		echo "  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
	fi

# Rollback database migrations
migrate-down:
	@echo "Rolling back migrations..."
	@if command -v migrate > /dev/null; then \
		migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down; \
	else \
		echo "Error: migrate tool not installed"; \
	fi

# Generate sqlc code
sqlc-gen:
	@echo "Generating sqlc code..."
	@if command -v sqlc > /dev/null; then \
		sqlc generate; \
		echo "sqlc generation complete"; \
	else \
		echo "Error: sqlc not installed. Install with:"; \
		echo "  go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest"; \
	fi

# Run linter
lint:
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "Error: golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Format complete"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies updated"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	@cp .env.example .env
	@make install-tools
	@make docker-up
	@sleep 5
	@make migrate-up
	@make sqlc-gen
	@echo "Setup complete! Run 'make run' to start the server"
