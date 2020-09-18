.DEFAULT_GOAL := build
.PHONY: build

build:
	go build -o $(GOPATH)/bin/anim8 cmd/anim8/main.go
