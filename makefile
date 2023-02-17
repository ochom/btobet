test:
	@echo "Testing ..."
	go test -v -coverprofile=coverage.out ./...

lint:
	@echo "Linting ..."
	@golangci-lint run -v

tidy:
	@echo "Tyding up ..."
	@go mod tidy