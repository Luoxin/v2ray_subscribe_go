FROM troian/golang-cross:v1.16.5
# https://hub.docker.com/r/troian/golang-cross
#"https://docker.mirrors.ustc.edu.cn/",
#"https://hub-mirror.c.163.com/",
#"https://reg-mirror.qiniu.com"

MAINTAINER luoxin <luoxin.ttt@gmail.com>
WORKDIR /home

COPY . /home/

# goreleaser version
ARG GORELEASER_VERSION=0.173.2
ARG GORELEASER_FILENAME=goreleaser_amd64.deb
# 安装 goreleaser
RUN  #!/bin/bash \
	if [ ! -f "$${GORELEASER_FILENAME}" ]; then \
		dpkg -i /home/goreleaser_amd64.deb \
	else \
		curl -L -o ./${GORELEASER_FILENAME} "https://cdn.jsdelivr.net/gh/goreleaser/goreleaser@releases/download/v${GORELEASER_VERSION}/${GORELEASER_FILENAME}" && dpkg -i ./${GORELEASER_FILENAME} \
	fi

WORKDIR /home/
#ENTRYPOINT ["goreleaser", "--skip-validate" ,"--skip-publish" ,"--snapshot" ,"--rm-dist", "--debug"]
ENTRYPOINT ["bash", "/home/gorelease.sh"]
