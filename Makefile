.DEFAULT_GOAL := build

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	go build ./...

generate:
	gomarkdoc -u -o README.md ./commercelayer/ && cd examples/full && terraform-docs .
