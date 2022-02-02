BIN := "./bin/rotator"
DOCKER_IMG="rotator:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/rotator

test:
	go test -race -count=100 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

run:
	docker-compose -f ./deployment/docker-compose.yaml -p rotator up --build

stop:
	docker-compose -f ./deployment/docker-compose.yaml -p rotator down

generate:
	go generate ./...

up-i:
	docker-compose -f ./deployment/docker-compose.yaml -p rotator up postgres rabbit

integration-tests:
	set -e ;\
	docker-compose -f ./deployment/docker-compose.test.yml up --build -d ;\
	test_status_code=0 ;\
	docker-compose -f ./deployment/docker-compose.test.yml run integration_tests go test -v || test_status_code=$$? ;\
	docker-compose -f ./deployment/docker-compose.test.yml down \
    --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
	exit $$test_status_code ;

.PHONY: build test install-lint-deps lint run stop generate up-i integration-tests
