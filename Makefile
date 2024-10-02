.PHONY: run test

run:
	@go run ./cmd/web

test:
	@go test -v ./... 
