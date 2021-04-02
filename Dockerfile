FROM dockercore/golang-cross

MAINTAINER luoxin <luoxin.ttt@gmail.com>
WORKDIR /build

ENV GOBIN=$GOPATH/bin
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY=https://goproxy.io,direct

RUN sed -i 's/deb.debian.org/mirrors.aliyun.com/g' /etc/apt/sources.list

COPY . /build/
#COPY ./build.sh /build/
#RUN chmod +x /build/build.sh
RUN rm -rf /build/go.sum

# 安装 goreleaser
RUN wget https://github.com/goreleaser/goreleaser/releases/download/v0.162.0/goreleaser_amd64.deb && dpkg -i goreleaser_amd64.deb
#RUN	dpkg -i /build/goreleaser_amd64.deb

# 安装golang
#ARG GOVERSION=1.16.2
#RUN wget https://studygolang.com/dl/go{{GOVERSION}}.linux-amd64.tar.gz && tar -zxvf go{{GOVERSION}}.linux-amd64.tar.gz && sudo mv go /usr/local/

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

#ENTRYPOINT ["build.sh"]
ENTRYPOINT ["goreleaser", "--skip-validate" ,"--skip-publish" ,"--snapshot" ,"--rm-dist"]
