.DEFAULT_GOAL := generate build

fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

test:
	go test -v ./...

build:
	go build ./...

generate:
	go generate