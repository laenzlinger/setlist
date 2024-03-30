.PHONY: help clean build run test
.DEFAULT_GOAL := help

DOCKER_IMAGE = ghcr.io/laenzlinger/setlist

RUN = docker run --rm -v $(shell pwd)/test/Repertoire:/repertoire $(DOCKER_IMAGE)

build: ## build the binary
	go build -o setlist main.go

run: ## run the application
	go run main.go

test: ## run tests
	go test ./...

integration-test: clean docker-build
	$(RUN) sheet --band Band --all
	$(RUN) sheet --band Band --gig "Grand Ole Opry"
	$(RUN) list --band Band --gig "Grand Ole Opry"

lint: ## lint source code
	golangci-lint run

clean: ## clean all output files
	rm -f setlist
	rm -rf test/Repertoire/out
	go clean -testcache

docker-build: build
	docker build -t $(DOCKER_IMAGE):latest .

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
