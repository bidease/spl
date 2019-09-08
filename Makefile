darwin:
	GOOS=darwin GOARCH=amd64 go build -o spl-amd64-darwin ./

linux:
	GOOS=linux GOARCH=amd64 go build -o spl-amd64-linux ./

.PHONY: darwin linux
