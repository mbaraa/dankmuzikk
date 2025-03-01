.PHONY: build

BINARY_NAME=dankmuzikk-server

build:
	npm i && \
	go mod tidy && \
	go generate && \
	go build ./cmd/server/main.go -ldflags="-w -s" -o ${BINARY_NAME}

generate:
	go generate

init:
	go mod tidy && \
	go get && \
    go run ./cmd/migrator/main.go

dev:
	go run github.com/cosmtrek/air@v1.51.0

beta: init
	GO_ENV="beta" ./${BINARY_NAME}

prod: init
	GO_ENV="prod" ./${BINARY_NAME}

clean:
	go clean

