VER := $(shell cat main.go | grep app.Version | cut -d " " -f 3)

build:
	go install ./

macos:
	GOOS=darwin GOARCH=amd64 go build -o spl ./

linux:
	GOOS=linux GOARCH=amd64 go build -o spl ./

release:
	GOOS=darwin GOARCH=amd64 go build -o spl ./ && tar zcf spl-macos-amd64-$(VER).tar.gz spl && rm spl
	GOOS=linux GOARCH=amd64 go build -o spl ./ && tar zcf spl-linux-amd64-$(VER).tar.gz spl  && rm spl

.PHONY: build macos linux release
