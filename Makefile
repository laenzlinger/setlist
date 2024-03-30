.PHONY: help clean build run test
.DEFAULT_GOAL := help

RUN = docker run --rm -v $(shell pwd)/test/Repertoire:/repertoire setlist

build: ## build the binary
	go build -o bin/setlist main.go

run: ## run the application
	go run main.go

test: ## run tests
	go test ./...

integration-test: clean docker-build
	$(RUN) setlist sheet --band Band --all
	$(RUN) setlist sheet --band Band --gig "Grand Ole Opry"
	$(RUN) setlist list --band Band --gig "Grand Ole Opry"

lint: ## lint source code
	golangci-lint run

clean: ## clean all output files
	rm -rf bin
	go clean -testcache

docker-build: build
	docker build -t setlist:latest .

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
