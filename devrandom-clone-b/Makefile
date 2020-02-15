BINARY_NAME := devrandom-clone
BINARY_OUTPUT_PATH := ./bin/$(BINARY_NAME)
SHELL := /bin/bash
PLATFORM := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)
GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin
GOCMD :=go

default: clean test build

lint:
	golangci-lint run ./...

clean:
	rm -rf ./bin

build:
	go fmt ./...
	DEP_BUILD_PLATFORMS=$(PLATFORM) DEP_BUILD_ARCHS=$(ARCH) $(GOCMD) build -o $(BINARY_OUTPUT_PATH) -v

test:
	$(GOCMD) test -v -race ./...

install: build
	cp ./bin $(GOBIN)

.PHONY: build test clean install lint