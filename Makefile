.DEFAULT_GOAL := generate build

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	go build ./...

generate:
	go generate