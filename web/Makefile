.PHONY: build

BINARY_NAME=dankmuzikk-web

build: htmx-init tailwindcss-init
	go mod tidy && \
	go generate && \
	go build -ldflags="-w -s" -o ${BINARY_NAME}

generate: tailwindcss-build
	go generate

init: htmx-init tailwindcss-init
	go get && \
	go generate

htmx-init:
	mkdir -p static/js/htmx && \
	wget https://unpkg.com/hyperscript.org@0.9.14/dist/_hyperscript.min.js -O static/js/htmx/hyperscript.min.js &&\
	wget https://unpkg.com/htmx-ext-json-enc@2.0.2/dist/json-enc.min.js -O static/js/htmx/json-enc.js &&\
	wget https://unpkg.com/htmx-ext-loading-states@2.0.1/dist/loading-states.min.js -O static/js/htmx/loading-states.js &&\
	wget https://unpkg.com/htmx.org@2.0.4/dist/htmx.min.js -O static/js/htmx/htmx.min.js

tailwindcss-init:
	mkdir -p static/css &&\
	npm i &&\
	npx tailwindcss@3.4.1 build -i static/css/style.css -o static/css/tailwind.css -m

tailwindcss-build:
	npx tailwindcss@3.4.1 build -i static/css/style.css -o static/css/tailwind.css -m

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

