BIN := "./bin/monitoring"
DOCKER_IMG="monitoring:1.0"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="1.0" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/monitoring

build_linux:
	@env GOOS=linux GOARCH=amd64 go build -a -v -o $(BIN) -installsuffix cgo ./cmd/monitoring

build_darwin:
	@env GOOS=darwin GOARCH=amd64 go build -a -v -o $(BIN) -installsuffix cgo ./cmd/monitoring

run_server: build
	$(BIN) grpc_server --config ./configs/config.json

run_client: build
	$(BIN) grpc_client

build_img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f deployments/docker/Dockerfile .

run_img: build-img
	docker run $(DOCKER_IMG)

test:
	go test -race -count 5 ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.37.0

lint: install-lint-deps
	golangci-lint run ./...

lint_ci: install-lint-deps
	golangci-lint run --config=./.golangci.yml

.PHONY: build build_linux build_darwin run_server build_img run_img test lint lint_ci
