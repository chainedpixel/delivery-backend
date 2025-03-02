.PHONY: test test-unit test-integration test-coverage

swagger:
	swag init -g cmd/main.go -o docs/swagger

GO=go
GOTEST=$(GO) test
PACKAGES=$(shell go list ./... | grep -v /vendor/)

test:
	@echo "Running all tests"
	$(GOTEST) -v -cover $(PACKAGES)

test-unit:
	@echo "Running unit tests"
	$(GOTEST) -v -short $(PACKAGES)

test-integration:
	@echo "Running integration tests"
	$(GOTEST) -v -run Integration $(PACKAGES)

test-coverage:
	@echo "Running coverage tests"
	$(GOTEST) -v -coverprofile=coverage.out $(PACKAGES)
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated in coverage.html"

test-bench:
	@echo "Running benchmark tests"
	$(GOTEST) -v -bench=. $(PACKAGES)