.PHONY: build build-alpine clean test help default

BIN_NAME=krontab

KRONTAB_VERSION := $(shell git describe --tags --dirty)
GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "+CHANGES" || true)
BUILD_DATE=$(shell date '+%Y-%m-%d-%H:%M:%S')
IMAGE_NAME := "jacobtomlinson/krontab"
BUILD_FLAGS="-X github.com/jacobtomlinson/krontab/version.GitCommit=${GIT_COMMIT}${GIT_DIRTY} -X github.com/jacobtomlinson/krontab/version.BuildDate=${BUILD_DATE} -X github.com/jacobtomlinson/krontab/version.Version=${KRONTAB_VERSION}"

default: help

help:
	@echo 'Management commands for krontab:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	@echo '    make build-all       Compile the project for all supported architectures with goreleaser.'
	@echo '    make test            Run tests on a compiled project.'
	@echo '    make deploy          Perform a deployment of the current tag with goreleaser.'
	@echo '    make release         Create a git tag and push to GitHub.'
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "building ${BIN_NAME} ${VERSION}"
	@echo "GOPATH=${GOPATH}"
	go build -ldflags ${BUILD_FLAGS} -o bin/${BIN_NAME}

export KRONTAB_VERSION
export GIT_COMMIT
export GIT_DIRTY
export BUILD_DATE
build-all:
	curl -sL https://git.io/goreleaser | bash -s -- --snapshot --skip-publish --rm-dist

export KRONTAB_VERSION
export GIT_COMMIT
export GIT_DIRTY
export BUILD_DATE
deploy:
	curl -sL https://git.io/goreleaser | bash

release:
	@read -p "Enter tag [v{major}.{minor}.{patch}]: " tag; \
	echo Tagging $$tag \
	git checkout master && git pull origin master --tags && git tag -a $$tag -m "$$tag" && git push origin $$tag

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

test:
	go test ./...

