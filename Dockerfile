# base image
FROM golang:1.18-alpine as base
WORKDIR /app

ENV GO111MODULE="on"
ENV GOOS="linux"
ENV GOARCH="amd64"
ENV CGO_ENABLED=0

RUN apk update \
    && apk add --no-cache ca-certificates git \
    && update-ca-certificates

# dev / air image
FROM base AS dev
WORKDIR /app/src

RUN go install github.com/cosmtrek/air@latest \
    && go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 8080
EXPOSE 2345

ENTRYPOINT ["air"]

# builder
FROM base AS builder
WORKDIR /app

COPY ./app/ /app/
RUN go mod download \
    && go mod verify

RUN go build -a -o .

# prod: todo

