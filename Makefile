.PHONY: all test clean

install-air:
	go install github.com/air-verse/air@latest

build: test
	go build -o bin/ui -ldflags="-w -s" internal/cmd/main.go

install-lint:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

lint: install-lint
	golangci-lint run

tidy:
	go mod tidy

test:
	go test -cover ./...

cover-report:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

test-verbose:
	go test -cover -v ./...

air: test install-air
	air run internal/cmd/main.go

clean:
	rm -rf bin/* tmp/* logs/*.log