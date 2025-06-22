.PHONY: test test-verbose test-coverage test-race clean build run

# Test commands
test:
	go test ./...

test-verbose:
	go test -v ./...

test-coverage:
	go test -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

test-race:
	go test -race ./...

# Build commands
build:
	go build -o bin/payment main.go

run:
	go run main.go

# Clean up
clean:
	rm -f bin/payment coverage.out coverage.html

# Install dependencies
deps:
	go mod download
	go mod tidy

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	golangci-lint run

# Run all checks
check: fmt test-race test-coverage lint 