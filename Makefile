.PHONY: run build test clean deps docker-up docker-down

BINARY_NAME=intelligent-pioneer
BINARY_PATH=./bin/$(BINARY_NAME)

# --- Application Commands ---
run:
	@echo "Running the application..."
	@go run ./cmd/intelligent-pioneer/main.go

build:
	@echo "Building the application..."
	@go build -o $(BINARY_PATH) ./cmd/intelligent-pioneer/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

clean:
	@echo "Cleaning up..."
	@go clean
	@rm -f $(BINARY_PATH)

# --- Go Module Commands ---
deps:
	@echo "Installing dependencies..."
	@go mod tidy

# --- Docker Commands ---
docker-up:
	@echo "Starting Docker services..."
	@docker compose up -d

docker-down:
	@echo "Stopping Docker services..."
	@docker compose down