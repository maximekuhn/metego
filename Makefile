.DEFAULT_GOAL := help

fmtCheck:
	gofmt -l -s .

fmt:
	gofmt -l -s -w .

test_integration:
	go test -tags=integration -v ./...

test:
	go test -v ./...

build: templ
	go build -v -o bin/web cmd/web/main.go

build_rpi: templ
	./build/build_rpi.sh

clean:
	go clean
	rm -rf ./bin

lint:
	golangci-lint run

templ:
	templ generate

help:
	@echo "Available targets"
	@echo "fmtCheck         - Check if code is correctly formatted"
	@echo "fmt              - Format code"
	@echo "test_integration - Run all integration tests"
	@echo "test             - Run all tests"
	@echo "build            - Build all binaries"
	@echo "build_rpi        - Build all binaries for Raspberry Pi (ARMv7) // REQUIRES DOCKER"
	@echo "clean            - Clean up bin/ directory"
	@echo "lint             - Run linter"
	@echo "templ            - Generate all templ components"
