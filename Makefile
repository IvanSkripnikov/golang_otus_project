DOCKER_COMPOSE?=docker-compose
RUN=$(DOCKER_COMPOSE) run --rm app

BIN := "./bin/banner-rotation"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	$(DOCKER_COMPOSE) pull --ignore-pull-failures
	$(DOCKER_COMPOSE) build --force-rm --pull

run: build
	$(BIN) -config ./.env.example

up: run
	$(DOCKER_COMPOSE) up -d --remove-orphans

stop:
	$(DOCKER_COMPOSE) stop

version: build
	$(BIN) version

test:
	go test -race ./app/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run up version test lint