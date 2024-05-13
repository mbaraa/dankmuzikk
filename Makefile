.PHONY: build

BINARY_NAME=dankmuzikk

# build builds the tailwind css sheet, and compiles the binary into a usable thing.
build:
	npm i && \
	go mod tidy && \
	templ generate && \
	go generate && \
	go build -ldflags="-w -s" -o ${BINARY_NAME}

generate:
	go generate && \
    templ generate

init:
	npm i && \
	go get && \
	go generate && \
	templ generate && \
    go run main.go migrate

seed:
	go run main.go seed

# dev runs the development server where it builds the tailwind css sheet,
# and compiles the project whenever a file is changed.
dev:
	templ generate --watch --cmd="./run.sh dev"

clean:
	go clean

