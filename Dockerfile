FROM golang:1.11.0-alpine
LABEL maintainer "Kalan Chen <kjj6198@gmail.com>"

COPY vendor /go/src/github.com/kjj6198/drink-bot/vendor

COPY utils /go/src/github.com/kjj6198/drink-bot/utils
COPY services /go/src/github.com/kjj6198/drink-bot/services
COPY apis /go/src/github.com/kjj6198/drink-bot/apis
COPY app /go/src/github.com/kjj6198/drink-bot/app
COPY config /go/src/github.com/kjj6198/drink-bot/config
COPY db /go/src/github.com/kjj6198/drink-bot/db
COPY middlewares /go/src/github.com/kjj6198/drink-bot/middlewares
COPY models /go/src/github.com/kjj6198/drink-bot/models
COPY main.go /go/src/github.com/kjj6198/drink-bot/main.go
WORKDIR /go/src/github.com/kjj6198/drink-bot
RUN GOOS=linux go build -o main .
ENTRYPOINT /go/src/github.com/kjj6198/drink-bot/main
EXPOSE 8080