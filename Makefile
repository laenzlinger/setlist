.PHONY: help clean build run test
.DEFAULT_GOAL := help

build: ## build the binary
	go build -o bin/setlist main.go

run: ## run the application
	go run main.go

test: ## run tests
	go test ./...

integration-test: clean build
	cd test/Repertoire && ../../bin/setlist sheet --band Band --all
	cd test/Repertoire && ../../bin/setlist sheet --band Band --gig "Grand Ole Opry"
	cd test/Repertoire && ../../bin/setlist list --band Band --gig "Grand Ole Opry"

lint: ## lint source codoe
	golangci-lint run

clean: ## clean all output files
	rm -rf bin
	go clean -testcache

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
