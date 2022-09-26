.DEFAULT_GOAL := generate build

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	go build ./...

generate: dependencies
	go generate

dependencies:
	go get -u github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs