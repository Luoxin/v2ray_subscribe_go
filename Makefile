PWD = $(shell pwd)

.PHONY: run

run:
	go mod download
	CGO_ENABLED=1 go run ./cmd/eutamias.go

build:
	cp config.yaml.simple config.yaml
	go mod download
	go mod vendor
	docker build -t sub:latest .
	docker run -it -v ${PWD}:/build/ sub:latest

build1:
	rm -rf ./go.sum
	rm -rf ./vendor
	#go mod download
#	go mod vendor
	docker build -t sub:latest .
	docker run -it -v E:/Eutamias:/build/ sub:latest
	rm -rf ./go.sum
	#rm -rf ./vendor

build2:
	rm -rf ./go.sum
	#rm -rf ./vendor
	#go mod download
	#go mod vendor
	docker build -t sub:latest .
	docker run -it -v D:/develop/Eutamias:/build/ sub:latest
	rm -rf ./go.sum
	rm -rf ./vendor

sync1:
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ' -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" -o D:/Service/Subscribe/eutamias.exe ./cmd/eutamias.go
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ'" -o D:/Service/Subscribe/checkwall.exe ./tool/checkwall/checkwall.go

sync2:
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ' -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" -o D:/Server/Subscribe/eutamias.exe ./cmd/eutamias.go
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ'" -o D:/Server/Subscribe/checkwall.exe ./tool/checkwall/checkwall.go
