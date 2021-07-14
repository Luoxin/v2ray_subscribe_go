PWD = $(shell pwd)
GIT_TAG = $(shell git describe --tags )

.PHONY: run

run:
	go mod tidy
	CGO_ENABLED=1 go run  -ldflags "-X 'geolite.GeoLiteUrl=https://kutt.luoxin.live/GHfTBv' -X 'proxies.ClashTplUrl=https://kutt.luoxin.live/dxvcRb'" ./cmd/eutamias.go

dockerrun:
	docker build --no-cache -t eutamias:${GIT_TAG} .
	docker run -it -p 2000:2000 eutamias:${GIT_TAG}

build:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v ${PWD}:/build/ eutamias-pkg:latest

build1:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v E:/Eutamias:/home/ eutamias-pkg:latest

build2:
	docker build -f ./internal/Dockerfile -t eutamias-pkg:latest .
	docker run -it --env-file ./.env -v D:/develop/Eutamias:/home/ eutamias-pkg:latest

syncl:
	goreleaser --skip-publish --skip-validate --rm-dist  --config .goreleaser-linux.yml --debug
	- cp ./dist/eutamias_linux_amd64/eutamias $(EUTAMIAS_HOME)
	- cp ./dist/checkwall_linux_amd64/checkwall $(EUTAMIAS_HOME)
	- cp ./dist/proxycheck_linux_amd64/proxycheck $(EUTAMIAS_HOME)
	- cp ./dist/dnsquery_linux_amd64/dnsquery $(EUTAMIAS_HOME)
	- cp ./dist/tohru_linux_amd64/tohru $(EUTAMIAS_HOME)

syncw:
	goreleaser.exe --skip-publish --skip-validate --rm-dist  --config .goreleaser-windows.yml --debug
	- cp ./dist/eutamias_windows_amd64/eutamias.exe $(EUTAMIAS_HOME)
	- cp ./dist/checkwall_windows_amd64/checkwall.exe $(EUTAMIAS_HOME)
	- cp ./dist/proxycheck_windows_amd64/proxycheck.exe $(EUTAMIAS_HOME)
	- cp ./dist/dnsquery_windows_amd64/dnsquery.exe $(EUTAMIAS_HOME)
	- cp ./dist/tohru_windows_amd64/tohru.exe $(EUTAMIAS_HOME)
