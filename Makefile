PWD = $(shell pwd)

.PHONY: build cclinux

buildLinux:
	go mod download
	go mod vendor
	docker build -f Dockerfile_linux -t linuxSub:latest .

runLinux:
	docker run -it -v $(PWD):/build/ --rm linuxSub:latest
