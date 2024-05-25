.PHONY: build

BINARY_NAME=dankmuzikk

# build builds the tailwind css sheet, and compiles the binary into a usable thing.
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

# dev runs the development server where it builds the tailwind css sheet,
# and compiles the project whenever a file is changed.
dev:
	go run github.com/cosmtrek/air@v1.51.0

clean:
	go clean

