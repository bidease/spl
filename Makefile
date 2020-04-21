darwin:
	GOOS=darwin GOARCH=amd64 go build -o spl-amd64-darwin ./cmd/spl

linux:
	GOOS=linux GOARCH=amd64 go build -o spl-amd64-linux ./cmd/spl

.PHONY: darwin linux
