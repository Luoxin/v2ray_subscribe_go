PWD = $(shell pwd)

.PHONY: build

build:
	cp config.yaml.simple config.yaml
	go mod download
	go mod vendor
	docker build -t sub:latest .
	docker run -it -v ${PWD}:/build/ sub:latest
