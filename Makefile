.DEFAULT_GOAL := build

fmt:
	go fmt ./...
.PHONY: fmt

vet: fmt
	go vet ./...
.PHONY: vet

build: vet
	go mod tidy
	go build -o ./cmd/server/bin/mahjong-linebot ./cmd/server
.PHONY: build
