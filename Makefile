.PHONY: lint build test cover check clean

lint:
	golangci-lint run ./...

build:
	go build ./...

test:
	go test -race -coverprofile=coverage.out -covermode=atomic ./...

cover: test
	go tool cover -func=coverage.out
	@echo "--- Uncovered functions (should be empty) ---"
	@go tool cover -func=coverage.out | grep -v "100.0%" | grep -v "^total" || true

check: lint build test

clean:
	rm -f coverage.out
