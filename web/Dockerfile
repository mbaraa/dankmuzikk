FROM golang:1.24-alpine AS build

WORKDIR /app
COPY . .

RUN go install github.com/a-h/templ/cmd/templ@v0.3.857 &&\
    apk add --no-cache make npm nodejs

RUN make

FROM alpine:latest AS run

WORKDIR /app
COPY --from=build /app/dankmuzikk-web ./dankmuzikk-web

EXPOSE 3003

CMD ["./dankmuzikk-web"]

