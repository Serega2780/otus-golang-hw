BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_sender"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	rm -rf api/pb
	mkdir -p api/pb

	protoc \
    	--proto_path=proto/ \
    	--go_out=./ \
    	--go-grpc_out=./ \
    	proto/events/*.proto

	protoc -I proto/ --grpc-gateway_out ./ \
        --grpc-gateway_opt generate_unbound_methods=true \
        proto/events/*.proto

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/calendar_scheduler
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/calendar_sender

# run: build
# 	$(BIN_CALENDAR) -config ./configs/calendar_config.yaml

# build-img:
# 	docker build \
# 		--build-arg=LDFLAGS="$(LDFLAGS)" \
# 		-t $(DOCKER_IMG) \
# 		-f build/Dockerfile .
#
# run-img: build-img
# 	docker run $(DOCKER_IMG)
#
# version: build
# 	$(BIN_CALENDAR) version

up:
	docker compose -f docker-compose.yaml up -d

down:
	docker compose -f docker-compose.yaml down

test:
	go test -race ./internal/storage/memory ./internal/server/grpc ./internal/server/http

integration_test:
	set -e ;\
	docker compose -f docker-compose-test.yaml build ;\
	test_status_code=0 ;\
	docker compose -f docker-compose-test.yaml run integration_tests go test ./integration || test_status_code=$$? ;\
	docker compose -f docker-compose-test.yaml down ;\
	exit $$test_status_code ;

integration_test-cleanup:
	docker compose -f docker-compose-test.yaml down \
        --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 60; \
  	docker compose rm -f


install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

# run build-img run-img

.PHONY: build  version test lint
