FROM golang:1.22-alpine as build

WORKDIR /app
COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest &&\
    apk add --no-cache make npm nodejs &&\
    make

FROM alpine:latest as run

RUN apk add --no-cache make

WORKDIR /app
COPY --from=build /app/dankmuzikk ./dankmuzikk
COPY --from=build /app/run.sh ./run.sh

EXPOSE 8080

CMD ["make", "prod"]
