include .env
export

# Build the application
all: build test

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	docker compose -f docker-compose.yml up --build -d

# Shutdown DB container
docker-down:
	docker compose -f docker-compose.yml down

# Run the application in development mode
docker-dev:
	@echo "Building and starting containers in development mode..."
	docker compose -f docker-compose.dev.yml up --build -d

# Shutdown development mode
docker-dev-down:
	@echo "Stopping development containers..."
	docker compose -f docker-compose.dev.yml down

migrate-up:
	@goose -dir ./internal/database/migrations/ postgres "$(DATABASE_URL)" up

migrate-down:
	@goose -dir ./internal/database/migrations/ postgres "$(DATABASE_URL)" down

create-migration:
	@read -p "Enter migration name: " name; \
	if [ -z "$$name" ]; then \
		echo "Migration name cannot be empty."; \
		exit 1; \
	fi; \
	goose -dir ./internal/database/migrations/ create "$$name" sql

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v
# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch docker-run docker-down itest
