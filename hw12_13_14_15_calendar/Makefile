BIN := "./bin/calendar"
SCHDEULER := "./bin/calendar_scheduler"
SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN)_scheduler -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN)_sender -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

generate:
	rm -rf internal/server/grpc/pb
	mkdir -p internal/server/grpc/pb
	protoc --proto_path=api/ --go_out=internal/server/grpc/pb --go-grpc_out=internal/server/grpc/pb api/*.proto

migrate:
	goose -env configs/.db.env up

run: build migrate
	$(BIN) -config ./configs/calendar_config.toml

scheduler-run: build
	$(SCHDEULER) -config ./configs/scheduler_config.toml

sender-run: build
	$(SENDER) -config ./configs/sender_config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -v -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
