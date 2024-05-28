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

seed:
	go run main.go seed

dev:
	go run github.com/cosmtrek/air@v1.51.0

beta:
	./${BINARY_NAME} serve

prod:
	./${BINARY_NAME} migrate
	./${BINARY_NAME} serve

clean:
	go clean

