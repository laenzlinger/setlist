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

test: ## run tests
	go test ./...

integration-test: clean docker-build
	$(RUN) generate sheet --band-name Band --all
	$(RUN) generate sheet --band-name Band "Grand Ole Opry"
	$(RUN) generate list  --band-name Band "Grand Ole Opry"
	ls -lart test/Repertoire/out

lint: ## lint source code
	golangci-lint run

clean: ## clean all output files
	rm -f setlist
	rm -rf dist
	rm -rf test/Repertoire/out
	rm -f test/Repertoire/Band/Songs/Frankie\ and\ Johnnie.pdf
	go clean -testcache

docker-build: build
	docker build -t $(DOCKER_IMAGE):latest .

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
