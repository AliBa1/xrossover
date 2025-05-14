.DEFAULT_GOAL := build

.PHONY:fmt vet build test
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	go build

test:
	go test ./...
