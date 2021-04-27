PWD = $(shell pwd)

.PHONY: build

build:
	cp config.yaml.simple config.yaml
	go mod download
	go mod vendor
	docker build -t sub:latest .
	docker run -it -v ${PWD}:/build/ sub:latest

build1:
	rm -rf ./go.sum
	rm -rf ./vendor
	go mod download
	go mod vendor
	docker build -t sub:latest .
	docker run -it -v E:/Eutamias:/build/ sub:latest
	rm -rf ./go.sum
	rm -rf ./vendor

build2:
	rm -rf ./go.sum
	rm -rf ./vendor
	go mod download
	go mod vendor
	docker build -t sub:latest .
	docker run -it -v D:/develop/Eutamias:/build/ sub:latest
	rm -rf ./go.sum
	rm -rf ./vendor

sync1:
	CGO_ENABLED=1 go build -o D:/Service/Subscribe/eutamias.exe ./cmd/eutamias.go

sync2:
	CGO_ENABLED=1 go build -o D:/Server/Subscribe/eutamias.exe ./cmd/eutamias.go
