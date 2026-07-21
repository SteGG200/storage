LDFLAGS:=-s -w
.PHONY: all test verify fmt lint build clean

all: fmt lint build

test:
	go test -race -v ./...

verify:
	golangci-lint config verify

fmt:
	gofmt -s -w .

lint:
	golangci-lint run ./...

build:
	go build -ldflags="${LDFLAGS}" -tags production -o bin/storage ./cmd/storage

clean:
	go clean -cache
	rm -rf bin
