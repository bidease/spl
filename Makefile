
build:
	go install ./

build-macos:
	GOOS=darwin GOARCH=amd64 go build -o spl ./

build-linux:
	GOOS=linux GOARCH=amd64 go build -o spl ./

.PHONY: build build-macos build-linux
