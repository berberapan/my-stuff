.PHONY: build run test clean

run:
	@go run ./cmd/web

test:
	@go test -v ./... 
