.PHONY: build

build:
	GOOS=linux GOARCH=amd64 go build -ldflags '-w -s' -tags 'netgo osusergo' -o dist/devproxy-linux-amd64 cmd/devproxy/main.go
	GOOS=darwin GOARCH=amd64 go build -ldflags '-w -s' -tags 'netgo osusergo' -o dist/devproxy-darwin-amd64 cmd/devproxy/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags '-w -s' -tags 'netgo osusergo' -o dist/devproxy-windows-amd64.exe cmd/devproxy/main.go
