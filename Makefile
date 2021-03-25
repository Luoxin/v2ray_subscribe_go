PWD = $(shell pwd)

.PHONY: build cclinux

build:
	docker build -t sub:latest .

cclinux:
	docker run -it -v $(PWD):/app/code --rm gosub:latest
