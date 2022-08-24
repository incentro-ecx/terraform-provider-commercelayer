.DEFAULT_GOAL := build

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	go build ./...

