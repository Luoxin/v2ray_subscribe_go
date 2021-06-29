PWD = $(shell pwd)

.PHONY: run

run:
	go mod download
	CGO_ENABLED=1 go run -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" ./cmd/eutamias.go

build:
	docker build -t eutamias:latest .
	docker run -it --env-file ./.env -v ${PWD}:/build/ eutamias:latest

build1:
	docker build -t eutamias:latest .
	docker run -it --env-file ./.env -v E:/Eutamias:/build/ eutamias:latest

build2:
	docker build -t eutamias:latest .
	docker run -it --env-file ./.env -v D:/develop/Eutamias:/build/ eutamias:latest

sync1:
	cp ./resource/clashTpl D:/Service/Subscribe/resource/
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ' -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" -o D:/Service/Subscribe/eutamias.exe ./cmd/eutamias.go
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ'" -o D:/Service/Subscribe/checkwall.exe ./tool/checkwall/checkwall.go

sync2:
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ' -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" -o D:/Server/Subscribe/eutamias.exe ./cmd/eutamias.go
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ'" -o D:/Server/Subscribe/checkwall.exe ./tool/checkwall/checkwall.go
