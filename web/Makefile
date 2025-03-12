.PHONY: build

BINARY_NAME=dankmuzikk-web

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
	go generate

dev:
	air -v > /dev/null
	@if [ $$? != 0 ]; then \
		echo "air was not found, installing it..."; \
		go install github.com/cosmtrek/air@v1.51.0; \
	fi
	air

beta:
	GO_ENV="beta" ./${BINARY_NAME}

prod:
	./${BINARY_NAME} migrate
	GO_ENV="prod" ./${BINARY_NAME}

clean:
	go clean

