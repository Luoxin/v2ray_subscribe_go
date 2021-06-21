FROM troian/golang-cross:v1.16.5
#https://hub.docker.com/r/troian/golang-cross
#"https://docker.mirrors.ustc.edu.cn/",
#"https://hub-mirror.c.163.com/",
#"https://reg-mirror.qiniu.com"

MAINTAINER luoxin <luoxin.ttt@gmail.com>
WORKDIR /build

ENV GOBIN=$GOPATH/bin
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY=https://goproxy.io,direct
ENV HOSTNAME=eutamias

COPY . /build/

# goreleaser version
ARG GORELEASER_VERSION=0.171.0
ARG GORELEASER_FILENAME=goreleaser_amd64.deb
# 安装 goreleaser
RUN  #!/bin/bash \
	if [ ! -f "$${GORELEASER_FILENAME}" ]; then \
		dpkg -i /build/goreleaser_amd64.deb \
	else \
		wget "https://cdn.jsdelivr.net/gh/goreleaser/goreleaser@releases/download/v${GORELEASER_VERSION}/${GORELEASER_FILENAME}" && dpkg -i ${GORELEASER_FILENAME} \
	fi

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list

RUN apt-get update && \
    apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    software-properties-common

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | apt-key add - && \
 	add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/debian \
   $(lsb_release -cs) \
   stable"

RUN apt-get update && \
	apt-get install -y docker-ce \
	docker-ce-cli

ENTRYPOINT ["goreleaser", "--skip-validate" ,"--skip-publish" ,"--snapshot" ,"--rm-dist", "--debug"]
