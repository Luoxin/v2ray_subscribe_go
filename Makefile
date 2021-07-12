PWD = $(shell pwd)

.PHONY: run

run:
	go mod download
	CGO_ENABLED=1 go run -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" ./cmd/eutamias.go

dockerrun:
	docker build --no-cache -t eutamias:0.0.3 .
	docker run -it -p 2000:2000 eutamias:0.0.3

build:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v ${PWD}:/build/ eutamias-pkg:latest

build1:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v E:/Eutamias:/home/ eutamias-pkg:latest

build2:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v D:/develop/Eutamias:/home/ eutamias-pkg:latest

sync1:
	goreleaser.exe --skip-publish --skip-validate --rm-dist  --config .goreleaser-windows.yml --debug
	cp ./dist/eutamias_windows_amd64/eutamias.exe D:/Service/Subscribe/
	cp ./dist/checkwall_windows_amd64/checkwall.exe D:/Service/Subscribe/
	cp ./dist/proxycheck_windows_amd64/proxycheck.exe D:/Service/Subscribe/
	rm -rf ./dist

sync2:
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' -X 'main.UpdateUrl=https://kutt.luoxin.live/0NnXIQ' -X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" -o D:/Server/Subscribe/eutamias.exe ./cmd/eutamias.go
	#CGO_ENABLED=1 go build -o D:/Server/Subscribe/eutamias.exe ./cmd/eutamias.go
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic' " -o D:/Server/Subscribe/checkwall.exe ./tool/checkwall/.
	CGO_ENABLED=1 go build -ldflags "-s -w --extldflags '-static -fpic'" -o D:/Server/Subscribe/proxycheck.exe ./tool/proxycheck/.
