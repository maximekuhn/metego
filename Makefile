.DEFAULT_GOAL := help

fmtCheck:
	gofmt -l -s .

fmt:
	gofmt -l -s -w .

test_integration:
	go test -tags=integration -v ./...

help:
	@echo "Available targets"
	@echo "fmtCheck         - Check if code is correctly formatted"
	@echo "fmt              - Format code"
	@echo "test_integration - Run all integration tests"
