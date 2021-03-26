FROM golang:1.16.2
MAINTAINER luoxin <luoxin.ttt@gmail.com>

WORKDIR /build

ENV GOBIN=$GOPATH/bin
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY=https://goproxy.io,direct

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list
RUN apt-get update && apt-get install -y --no-install-recommends gcc

COPY ./. /build/.
RUN rm -rf /build/go.sum
RUN chmod +x /build/docker-entrypoint.sh

ENTRYPOINT ["bash", "./docker-entrypoint.sh"]
