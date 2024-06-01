.DEFAULT_GOAL := help

fmtCheck:
	gofmt -l -s .

fmt:
	gofmt -l -s -w .

test_integration:
	go test -tags=integration -v ./...

build:
	go build -v -o bin/cli cmd/cli/main.go

clean:
	go clean
	rm -rf ./bin

lint:
	golangci-lint run

help:
	@echo "Available targets"
	@echo "fmtCheck         - Check if code is correctly formatted"
	@echo "fmt              - Format code"
	@echo "test_integration - Run all integration tests"
	@echo "build            - Build all binaries"
	@echo "clean            - Clean up bin/ directory"
	@echo "lint             - Run linter"
