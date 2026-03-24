.PHONY: dev build test clean

dev:
	air

build:
	go build -o bin/sweatshop ./cmd/server

test:
	go test -v ./...

clean:
	rm -rf bin/
