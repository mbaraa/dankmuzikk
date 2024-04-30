FROM golang:1.22-alpine as build

WORKDIR /app
COPY . .

RUN go install github.com/a-h/templ/cmd/templ@latest &&\
    apk add make npm nodejs &&\
    make

FROM alpine:latest as run

RUN apk add yt-dlp

WORKDIR /app
COPY --from=build /app/dankmuzikk ./run

EXPOSE 8080

CMD ["./run", "serve"]
