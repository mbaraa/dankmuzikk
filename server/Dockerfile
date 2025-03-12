FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.2.731 &&\
    apk add --no-cache make &&\
    make init &&\
    make build-server &&\
    make build-migrator

FROM alpine:latest AS run

RUN apk add --no-cache make

WORKDIR /app
COPY --from=build /app/dankmuzikk-server ./dankmuzikk-server
COPY --from=build /app/Makefile ./Makefile

EXPOSE 8080

CMD ["make", "dankmuzikk-server-prod"]
