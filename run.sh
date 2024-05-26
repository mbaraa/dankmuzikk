#!/bin/sh

beta() {
    ./dankmuzikk serve
}

prod() {
    ./dankmuzikk migrate
    ./dankmuzikk serve
}

if [ $1 == "dev" ]; then
    dev
fi

if [ $1 == "beta" ]; then
    beta
fi

if [ $1 == "prod" ]; then
    prod
fi
