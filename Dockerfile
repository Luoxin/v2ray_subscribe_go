FROM golang:alpine
MAINTAINER luoxin <luoxin.ttt@gmail.com>

WORKDIR /app

ENV GOBIN=$GOPATH/bin
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io,direct

#RUN apt-get update && apt-get install -y --no-install-recommends \
#    ca-certificates \
#    sudo \
#    git \
#    gcc \
#    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
#    && rm -rf /var/lib/apt/lists/*g

COPY go.mod go.sum ./
RUN go mod tidy
COPY . .

RUN go build -o main ./cmd/windows.go

ENV TZ=Asia/Shanghai

ENTRYPOINT ["go", "build"]

