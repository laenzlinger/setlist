.PHONY: help clean build run test
.DEFAULT_GOAL := help

build: ## build the binary
	go build -o bin/setlist main.go

run: ## run the application
	go run main.go

test: ## run tests
	go test ./...


clean: ## clean all output files
	rm -rf bin

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
