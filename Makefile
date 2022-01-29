BIN := "./bin/rotator"
DOCKER_IMG="rotator:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/rotator

run: build
	$(BIN) -config ./configs/config.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--build-arg=CONFIG_FILE_NAME=config \
		-t $(DOCKER_IMG) \
		-f ./deployment/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

test:
	go test -race -count=1 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

up:
	docker-compose -f ./deployment/docker-compose.yaml -p rotator up

down:
	docker-compose -f ./deployment/docker-compose.yaml -p rotator down

rebuild:
	docker-compose -f ./deployment/docker-compose.yaml -p rotator up --build

up-i:
	docker-compose -f ./deployment/docker-compose.yaml -p rotator up postgres rabbit

generate:
	mkdir -p internal/server/bannersrotatorpb
	protoc -I ./api \
    --go_out ./internal/server/bannersrotatorpb/ --go_opt paths=source_relative \
    --go-grpc_out ./internal/server/bannersrotatorpb/ --go-grpc_opt paths=source_relative \
    api/*.proto
