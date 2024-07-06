.PHONY: build

BINARY_NAME=dankmuzikk

build:
	npm i && \
	go mod tidy && \
	go generate && \
	go build -ldflags="-w -s" -o ${BINARY_NAME}

generate:
	go generate

init:
	npm i && \
	go get && \
	go generate && \
    go run main.go migrate

dev:
	go run github.com/cosmtrek/air@v1.51.0

beta:
	GO_ENV="beta" ./${BINARY_NAME} serve

prod:
	./${BINARY_NAME} migrate
	GO_ENV="prod" ./${BINARY_NAME} serve

clean:
	go clean

