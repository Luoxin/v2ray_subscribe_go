#FROM troian/golang-cross:v1.16.5
FROM golang:1.16.5

MAINTAINER luoxin <luoxin.ttt@gmail.com>
WORKDIR /home

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

RUN git clone https://hub.fastgit.org/Luoxin/Eutamias.git /home/eutamias/

WORKDIR /home/eutamias
RUN go mod tidy
RUN curl -L -o /home/eutamias/resource/GeoLite2.mmdb $GEOLITE2_URL
RUN curl -L -o /home/eutamias/resource/clashTpl $CLASHTPL_URL
RUN curl -L -o /home/eutamias/.eutamias.es $BASE_DB_URL
RUN curl -L -o /home/config.yaml $BASE_CONFIG_URL
RUN go build -o /home/eutamias/eutamias -v ./cmd/.
RUN rm -rf /home/eutamias/eutamias

ENTRYPOINT ["go", "run", "-v", "./cmd/.", "-c", "/home/config.yaml"]
