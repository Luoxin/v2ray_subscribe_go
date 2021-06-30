PWD = $(shell pwd)

.PHONY: run

run:
	go mod download
	CGO_ENABLED=1 go run -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" ./cmd/eutamias.go

dockerrun:
	docker build --no-cache -t eutamias:latest .
	docker run -it -v config.yaml:/home/config.yaml -p 1900:8080 eutamias:latest

build:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v ${PWD}:/build/ eutamias-pkg:latest

build1:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v E:/Eutamias:/build/ eutamias-pkg:latest

build2:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v D:/develop/Eutamias:/build/ eutamias-pkg:latest

sync1:
	cp ./resource/clashTpl D:/Service/Subscribe/resource/
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ' -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" -o D:/Service/Subscribe/eutamias.exe ./cmd/eutamias.go
	#CGO_ENABLED=1 go build -o D:/Service/Subscribe/eutamias.exe ./cmd/eutamias.go
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ'" -o D:/Service/Subscribe/checkwall.exe ./tool/checkwall/.
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic'" -o D:/Service/Subscribe/proxycheck.exe ./tool/proxycheck/.

sync2:
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ' -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" -o D:/Server/Subscribe/eutamias.exe ./cmd/eutamias.go
	#CGO_ENABLED=1 go build -o D:/Server/Subscribe/eutamias.exe ./cmd/eutamias.go
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ'" -o D:/Server/Subscribe/checkwall.exe ./tool/checkwall/.
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic'" -o D:/Server/Subscribe/proxycheck.exe ./tool/proxycheck/.
