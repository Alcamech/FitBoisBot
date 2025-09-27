# FitBoisBot Makefile

.PHONY: build run serve dev install uninstall db-up db-down db-login db-seed fmt help

# Build the application
build:
	@echo "Building fitboisbot..."
	@go build -ldflags "-X 'github.com/Alcamech/FitBoisBot/internal/version.GitCommit=$$(git rev-parse HEAD)' -X 'github.com/Alcamech/FitBoisBot/internal/version.BuildTime=$$(date)'" -o bin/fitboisbot ./cmd/fitboisbot

# Run the application (requires TOKEN environment variable)
run: build
	@echo "Running fitboisbot serve..."
	@./bin/fitboisbot serve

# Alias for run
serve: run

# Run with live reload during development
dev:
	@echo "Starting development server..."
	@go run ./cmd/fitboisbot/main.go serve

# Install to ~/.local/bin
install: build
	@echo "Installing fitboisbot to ~/.local/bin..."
	@mkdir -p ~/.local/bin
	@cp bin/fitboisbot ~/.local/bin/fitboisbot
	@chmod +x ~/.local/bin/fitboisbot
	@echo "Installed! Make sure ~/.local/bin is in your PATH"
	@echo "Usage: fitboisbot serve, fitboisbot announce \"message\""

# Uninstall from ~/.local/bin
uninstall:
	@echo "Removing fitboisbot from ~/.local/bin..."
	@rm -f ~/.local/bin/fitboisbot
	@echo "Uninstalled!"

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -cover ./...

# Start database
db-up:
	@echo "Starting database..."
	@docker compose up fitboisbot-db -d

# Stop database
db-down:
	@echo "Stopping database..."
	@docker compose down

# Login to database
db-login:
	@echo "Connecting to database..."
	@docker compose exec fitboisbot-db mariadb -u root -pfitboi_4er! fitbois

# Seed database with mock data
db-seed:
	@echo "Seeding database with mock data..."
	@docker compose exec -T fitboisbot-db mariadb -u fitboi_user -pfitboi_4er! fitbois < scripts/insert_mock_data.sql
	@echo "Mock data inserted successfully!"

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go mod download

# Show help
help:
	@echo "Available commands:"
	@echo "  build        - Build the application"
	@echo "  run/serve    - Build and run the bot server"
	@echo "  dev          - Run with go run for development"
	@echo "  install      - Install to ~/.local/bin"
	@echo "  uninstall    - Remove from ~/.local/bin"
	@echo "  test         - Run tests"
	@echo "  test-coverage- Run tests with coverage report"
	@echo "  db-up        - Start database container"
	@echo "  db-down      - Stop database container"
	@echo "  db-login     - Login to database container"
	@echo "  db-seed      - Seed database with mock test data"
	@echo "  fmt          - Format code"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install/update dependencies"
	@echo "  help         - Show this help"
	@echo ""
	@echo "CLI Commands:"
	@echo "  ./bin/fitboisbot serve                    - Start bot server"
	@echo "  ./bin/fitboisbot announce \"Message\"       - Send to all groups"
	@echo "  ./bin/fitboisbot announce --file msg.md   - Send from file"
	@echo "  ./bin/fitboisbot version                  - Show version info"
	@echo "  ./bin/fitboisbot help                     - Show CLI help"
	@echo ""
	@echo "Environment variables:"
	@echo "  TOKEN - Required Telegram bot token"
	@echo ""
	@echo "Quick start:"
	@echo "  1. cp .env.example .env"
	@echo "  2. Edit .env with your bot token"
	@echo "  3. make db-up"
	@echo "  4. make db-seed"
	@echo "  5. make dev"

# Default target
.DEFAULT_GOAL := help
