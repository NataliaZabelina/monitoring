BIN := "./bin/monitoring"
DOCKER_IMG="monitoring:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/monitoring

build_ci:
	@env GOOS=linux GOARCH=amd64 go build -a -v -o $(BIN) -installsuffix cgo ./cmd/monitoring

run: build
	$(BIN) -config ./configs/config.json

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f deployments/docker/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race -count 5 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

lint_ci: install-lint-deps
	golangci-lint run --config=./.golangci.yml

.PHONY: build run build-img run-img version test lint
