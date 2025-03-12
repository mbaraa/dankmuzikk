.PHONY: build

SERVER_BINARY_NAME=dankmuzikk-server
CDN_BINARY_NAME=dankmuzikk-cdn
EVENTHUB_BINARY_NAME=dankmuzikk-eventhub
MIGRATOR_BINARY_NAME=dankmuzikk-migrator

build: init generate build-server build-migrator build-eventhub build-cdn

build-server:
	go build -ldflags="-w -s" -o ${SERVER_BINARY_NAME} ./cmd/http-server/main.go

build-eventhub:
	go build -ldflags="-w -s" -o ${EVENTHUB_BINARY_NAME} ./cmd/eventhub/main.go

build-migrator:
	go build -ldflags="-w -s" -o ${MIGRATOR_BINARY_NAME} ./cmd/migrator/main.go

build-cdn:
	go build -ldflags="-w -s" -o ${CDN_BINARY_NAME} ./cmd/cdn/main.go

init:
	go mod tidy && \
	go generate cmd/http-server/main.go

migrate: build-migrator
	./${MIGRATOR_BINARY_NAME}

dev:
	air -v > /dev/null
	@if [ $$? != 0 ]; then \
		echo "air was not found, installing it..."; \
		go install github.com/cosmtrek/air@v1.51.0; \
	fi

	air

dankmuzikk-server-beta: migrate
	GO_ENV="beta" ./${SERVER_BINARY_NAME}

dankmuzikk-server-prod: migrate
	GO_ENV="prod" ./${SERVER_BINARY_NAME}

generate:
	go generate cmd/http-server/main.go

clean:
	go clean
