#FROM troian/golang-cross:v1.16.5
FROM golang:1.16.5

MAINTAINER luoxin <luoxin.ttt@gmail.com>
WORKDIR /home

ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=off

RUN git clone https://hub.fastgit.org/Luoxin/Eutamias.git /home/eutamias/

WORKDIR /home/eutamias
RUN go mod tidy
RUN curl -L -o /home/eutamias/resource/GeoLite2.mmdb https://kutt.luoxin.live/GHfTBv
RUN curl -L -o /home/eutamias/resource/clashTpl https://kutt.luoxin.live/dxvcRb
RUN curl -L -o /home/eutamias/.eutamias.es https://kutt.luoxin.live/EiFhJq
RUN curl -L -o /home/config.yaml https://kutt.luoxin.live/3v4DWp
RUN go build -o /home/eutamias/eutamias -v ./cmd/.
RUN rm -rf /home/eutamias/eutamias

ENTRYPOINT ["go", "run", "-v", "./cmd/.", "-c", "/home/config.yaml"]
