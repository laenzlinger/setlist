.PHONY: help clean build run test
.DEFAULT_GOAL := help

DOCKER_IMAGE = ghcr.io/laenzlinger/setlist

RUN ?= docker run --rm --user "$(shell id -u)":"$(shell id -g)" -v $(shell pwd)/test/Repertoire:/repertoire $(DOCKER_IMAGE)

build: ## build the binary
	go build -o setlist main.go

run: ## run the application
	go run main.go

install: ## run the application
	go install main.go

generate: ## generate code
	go generate ./...

test: ## run tests
	go test ./...

test-integration: clean docker-build ## run integration tests
	$(RUN) clean
	$(RUN) generate sheet --all
	$(RUN) generate sheet
	$(RUN) generate list --landscape --font-size 40px
	$(RUN) generate suisa "Grand Ole Opry"
	ls -lartR test/Repertoire/out

test-all: lint test test-integration

lint: ## lint source code
	golangci-lint run

clean: ## clean all output files
	rm -f setlist
	rm -rf dist
	go clean -testcache

docker-build: build ## build the ocker container
	docker build -t $(DOCKER_IMAGE):latest .

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
