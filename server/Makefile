.PHONY: build

SERVER_BINARY_NAME=dankmuzikk-server
CDN_BINARY_NAME=dankmuzikk-cdn
EVENTHUB_BINARY_NAME=dankmuzikk-eventhub
MIGRATOR_BINARY_NAME=dankmuzikk-migrator

all: build-server build-eventhub build-cdn build-migrator

build: init generate build-server build-migrator build-eventhub build-cdn

build-server: init generate
	go build -ldflags="-w -s" -o ${SERVER_BINARY_NAME} ./cmd/http-server/main.go

build-eventhub: init generate
	go build -ldflags="-w -s" -o ${EVENTHUB_BINARY_NAME} ./cmd/eventhub/...

build-migrator: init
	go build -ldflags="-w -s" -o ${MIGRATOR_BINARY_NAME} ./cmd/migrator/main.go

build-cdn: init
	go build -ldflags="-w -s" -o ${CDN_BINARY_NAME} ./cmd/cdn/main.go

init:
	go mod tidy && \
	go install github.com/a-h/templ/cmd/templ@v0.3.857

migrate: build-migrator
	./${MIGRATOR_BINARY_NAME}

dev:
	air -v > /dev/null
	@if [ $$? != 0 ]; then \
		echo "air was not found, installing it..."; \
		go install github.com/cosmtrek/air@v1.51.0; \
	fi

	air

dankmuzikk-server-beta:
	GO_ENV="beta" ./${SERVER_BINARY_NAME}

dankmuzikk-server-prod:
	./${MIGRATOR_BINARY_NAME}
	GO_ENV="prod" ./${SERVER_BINARY_NAME}

generate:
	templ generate -path ./mailer/

clean:
	go clean
