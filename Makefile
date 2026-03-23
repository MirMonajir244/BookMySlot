.PHONY: build run test clean docker-up docker-down docker-rebuild migrate help

# Variables
APP_NAME=bookmyslot
MAIN_PATH=./cmd/server
DOCKER_COMPOSE=docker compose

## help: Show this help message
help:
	@echo "BookMySlot - Available Commands:"
	@echo ""
	@echo "  make build          Build the Go binary"
	@echo "  make run            Run the server locally"
	@echo "  make test           Run all unit tests"
	@echo "  make test-verbose   Run all unit tests with verbose output"
	@echo "  make clean          Remove build artifacts"
	@echo "  make docker-up      Start all services (PostgreSQL + App)"
	@echo "  make docker-down    Stop all services"
	@echo "  make docker-rebuild Rebuild and start all services"
	@echo "  make docker-logs    View running container logs"
	@echo "  make lint           Run go vet"
	@echo "  make fmt            Format all Go files"
	@echo "  make tidy           Run go mod tidy"
	@echo ""

## build: Build the Go binary
build:
	@echo "🔨 Building $(APP_NAME)..."
	go build -o $(APP_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: ./$(APP_NAME)"

## run: Run the server locally (requires PostgreSQL running)
run:
	@echo "🚀 Starting $(APP_NAME)..."
	go run $(MAIN_PATH)/main.go

## test: Run all unit tests
test:
	@echo "🧪 Running tests..."
	go test ./tests/... -count=1
	@echo "✅ All tests passed"

## test-verbose: Run all unit tests with verbose output
test-verbose:
	@echo "🧪 Running tests (verbose)..."
	go test ./tests/... -v -count=1

## clean: Remove build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -f $(APP_NAME)
	@echo "✅ Clean complete"

## docker-up: Start all services with Docker Compose
docker-up:
	@echo "🐳 Starting services..."
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Services started at http://localhost:8080"

## docker-down: Stop all services
docker-down:
	@echo "🐳 Stopping services..."
	$(DOCKER_COMPOSE) down
	@echo "✅ Services stopped"

## docker-rebuild: Rebuild and start all services
docker-rebuild:
	@echo "🐳 Rebuilding and starting services..."
	$(DOCKER_COMPOSE) up --build -d
	@echo "✅ Services rebuilt and started at http://localhost:8080"

## docker-logs: View container logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

## lint: Run go vet
lint:
	@echo "🔍 Running go vet..."
	go vet ./...
	@echo "✅ No issues found"

## fmt: Format all Go files
fmt:
	@echo "✨ Formatting code..."
	go fmt ./...
	@echo "✅ Code formatted"

## tidy: Run go mod tidy
tidy:
	@echo "📦 Tidying dependencies..."
	go mod tidy
	@echo "✅ Dependencies tidied"
