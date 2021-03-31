PWD = $(shell pwd)

.PHONY: build cclinux

build:
	cp config.yaml.simple config.yaml
	go mod download
	go mod vendor
	docker build -t sub:latest .
	docker run -it -v ${PWD}:/build/ --rm sub:latest
