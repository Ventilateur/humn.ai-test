SHELL := /bin/bash

BUILD_DIR := build

test:
	mkdir -p $(BUILD_DIR)
	go test ./... -v -cover -coverprofile=$(BUILD_DIR)/coverage.out -coverpkg=./...
	go tool cover -func $(BUILD_DIR)/coverage.out

.PHONY: build
build:
	go build -o app
