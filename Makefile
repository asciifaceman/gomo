.PHONY: clean test build build-linux

VERSION=`cat VERSION`
PACKAGE:=github.com/asciifaceman/gomo

clean: ## clean build dir
	rm -rf target/

test: ## run tests
	go test ./... -cover

build: clean build-linux

build-linux:
	@GOOS=linux GOARCH=amd64 go build \
	-ldflags "-X ${PACKAGE}/cmd.version=${VERSION} -X '${PACKAGE}/cmd.build=$(shell date)'" \
	-o target/gomo
