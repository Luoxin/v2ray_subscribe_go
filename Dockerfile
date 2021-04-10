FROM troian/golang-cross:latest
#https://hub.docker.com/r/troian/golang-cross

MAINTAINER luoxin <luoxin.ttt@gmail.com>
WORKDIR /build

ENV GOBIN=$GOPATH/bin
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV GOPROXY=https://goproxy.io,direct

# 安装golang
# https://github.com/letseeqiji/oneinstall/blob/master/golang/goinstall.sh
#RUN wget https://studygolang.com/dl/go1.16.2.linux-amd64.tar.gz && tar -zxvf go1.16.2.linux-amd64.tar.gz && sudo mv go /usr/local/
# gourl=$(curl -s  https://studygolang.com/dl |  sed -n '/dl\/golang\/go.*\.linux-amd64\.tar\.gz/p' | sed -n '1p' | sed -n '/1/p' | awk 'BEGIN{FS="\""}{print $4}')
#RUN wget "https://studygolang.com/dl/golang/go1.16.3.linux-amd64.tar.gz"   &&  tar -C /usr/local -zxvf go1.16.3.linux-amd64.tar.gz

COPY . /build/
#COPY ./build.sh /build/
#RUN chmod +x /build/build.sh
RUN rm -rf /build/go.sum

# goreleaser version
ARG GORELEASER_VERSION=0.162.0
# 安装 goreleaser
RUN wget "https://github.com/goreleaser/goreleaser/releases/download/v${GORELEASER_VERSION}/goreleaser_amd64.deb" && dpkg -i goreleaser_amd64.deb
#RUN	dpkg -i /build/goreleaser_amd64.deb

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

#ENTRYPOINT ["build.sh"]
ENTRYPOINT ["goreleaser", "--skip-validate" ,"--skip-publish", "--debug" ,"--snapshot" ,"--rm-dist"]
