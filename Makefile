PWD = $(shell pwd)

.PHONY: build cclinux

build:
	go mod download
	go mod vendor
	docker build -t sub:latest .

cclinux:
	docker run -it -v $(PWD):/build/ --rm sub:latest
