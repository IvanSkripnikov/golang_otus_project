DOCKER_COMPOSE?=docker-compose

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

run:
	$(DOCKER_COMPOSE) up -d --build

down:
	$(DOCKER_COMPOSE) down

# запуск интеграционных тестов
test-integration: run
	cd app && \
	go test -race -count 1 .

# запуск тестов core-логики
test-core: run
	cd app && \
	cd components && \
	go test -race -count 100 .

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	cd app && \
	golangci-lint run ./...

.PHONY: build run version test lint