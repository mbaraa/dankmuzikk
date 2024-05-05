#!/bin/sh

dev() {
    # This function is used to compile tailwind's style sheet, and run the web server.
    # Also this is supposed to be ran using templ's CLI, i.e.
    # templ generate --watch --cmd="./run.sh"
    # idk, that's it.

    go generate
    go run . serve
}

prod() {
    ./dankmuzikk migrate
    ./dankmuzikk serve
}

if [ $1 == "dev" ]; then
    dev
fi

if [ $1 == "prod" ]; then
    prod
fi
