BIN := "./bin/calendar"
SCHDEULER := "./bin/calendar_scheduler"
SENDER := "./bin/calendar_sender"
DOCKER_CALENDAR_IMG="calendar"
DOCKER_SCHEDULER_IMG="scheduler"
DOCKER_SENDER_IMG="sender"
DOCKER_TESTS_IMG="integration-tests"

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

run: build
	$(BIN) -config ./configs/calendar_config.toml

scheduler-run: build
	$(SCHDEULER) -config ./configs/scheduler_config.toml

sender-run: build
	$(SENDER) -config ./configs/sender_config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_CALENDAR_IMG) \
		-f build/calendar/Dockerfile .

	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SCHEDULER_IMG) \
		-f build/scheduler/Dockerfile .

	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_SENDER_IMG) \
		-f build/sender/Dockerfile .
	
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_TESTS_IMG) \
		-f build/tests/Dockerfile .

run-img: build-img
	docker run -d $(DOCKER_CALENDAR_IMG)
	docker run -d $(DOCKER_SCHEDULER_IMG)
	docker run -d $(DOCKER_SENDER_IMG)

up:
	docker-compose -f deployments/docker-compose.yaml up

down:
	docker-compose -f deployments/docker-compose.yaml down

version: build
	$(BIN) version

test:
	go test -v -race ./internal/...

integration-tests:
	docker-compose -f deployments/docker-compose.test.yaml up --build --abort-on-container-exit --exit-code-from integration-tests
	docker-compose -f deployments/docker-compose.test.yaml down	

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.62.0

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint
