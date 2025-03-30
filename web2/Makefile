.PHONY: build

BINARY_NAME=dankmuzikk-web

build:
	go mod tidy && \
	go generate && \
	go build -ldflags="-w -s" -o ${BINARY_NAME}

generate: tailwindcss-build
	go generate

init: tailwindcss-init
	mkdir -p static/css &&\
	mkdir -p static/js/htmx &&\
	wget https://unpkg.com/htmx-ext-loading-states@2.0.0/loading-states.js -O static/js/htmx/loading-states.js &&\
	wget https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -O static/js/htmx/htmx.min.js &&\
	go get && \
	go generate

tailwindcss-init:
	npm i &&\
	npx @tailwindcss/cli -i static/css/style.css -o static/css/tailwind.css -m

tailwindcss-build:
	npx @tailwindcss/cli -i static/css/style.css -o static/css/tailwind.css -m

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
	GO_ENV="prod" ./${BINARY_NAME}

clean:
	go clean

