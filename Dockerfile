FROM golang:1.16.2
MAINTAINER luoxin <luoxin.ttt@gmail.com>

WORKDIR /app

ENV GOBIN=$GOPATH/bin
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY=https://goproxy.io,direct

RUN apt-get update && apt-get install -y --no-install-recommends gcc

COPY . .

RUN go build -o windows_docker ./cmd/windows.go
