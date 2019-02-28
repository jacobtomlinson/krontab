.PHONY: build build-alpine clean test help default

BIN_NAME=krontab

VERSION := $(shell git describe --tags --dirty)
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "jacobtomlinson/krontab"
BUILD_FLAGS="-X github.com/jacobtomlinson/krontab/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/jacobtomlinson/krontab/version.BuildDate=${BUILD_DATE} -X github.com/jacobtomlinson/krontab/version.Version=${VERSION}"

default: help

help:
	@echo 'Management commands for krontab:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make build-all       Compile the project for all supported architectures.'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags ${BUILD_FLAGS} -o bin/${BIN_NAME}

build-linux-amd64:
	@echo "building ${BIN_NAME} ${VERSION} - linux amd64"
	@echo "GOPATH=${GOPATH}"
	GOOS=linux GOARCH=amd64 go build -ldflags ${BUILD_FLAGS} -o bin/${BIN_NAME}-linux-amd64

build-linux-arm:
	@echo "building ${BIN_NAME} ${VERSION} - linux arm"
	@echo "GOPATH=${GOPATH}"
	GOOS=linux GOARCH=arm go build -ldflags ${BUILD_FLAGS} -o bin/${BIN_NAME}-linux-arm

build-darwin-amd64:
	@echo "building ${BIN_NAME} ${VERSION} - darwin arm"
	@echo "GOPATH=${GOPATH}"
	GOOS=darwin GOARCH=amd64 go build -ldflags ${BUILD_FLAGS} -o bin/${BIN_NAME}-darwin-amd64

build-all: build-linux-amd64 build-linux-arm build-darwin-amd64

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test ./...

