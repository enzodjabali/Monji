BINARY=bin/server
MAIN_SERVER=./cmd/api/main.go

build:
	@echo "Building the application..."
	@mkdir -p bin
	go build -o $(BINARY) $(MAIN_SERVER)

run:
	@echo "Running the application..."
	go run $(MAIN_SERVER)

test:
	@echo "Running tests..."
	go test ./... -v

coverage:
	@echo "Generating test coverage report..."
	go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

fmt:
	@echo "Formatting Go code..."
	go fmt ./...

tidy:
	@echo "Tidying go modules..."
	go mod tidy

clean:
	@echo "Cleaning up..."
	rm -rf bin coverage.out coverage.html

.PHONY: build run test coverage fmt tidy clean
