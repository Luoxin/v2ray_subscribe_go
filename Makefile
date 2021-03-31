PWD = $(shell pwd)

.PHONY: build cclinux

build:
	go mod download
	go mod vendor
	docker build -t sub:latest .
	docker run -it -v ${PWD}:/build/ --rm sub:latest
