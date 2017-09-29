SHELL = /bin/bash

TARGET       = stone
PROJECT_NAME = github.com/zanecloud/stone

MAJOR_VERSION = $(shell cat VERSION)
GIT_VERSION   = $(shell git log -1 --pretty=format:%h)

BUILD_IMAGE     = golang:1.8.3-onbuild

stone:clean
	CGO_ENABLED=0  go build -a -v -o ${TARGET}

docker-disk:clean-docker-disk
	cd docker-disk && CGO_ENABLED=0  go build -a -v -o docker-disk && cd ..

build:
	docker run --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} make stone
	docker run --rm -v $(shell pwd):/go/src/${PROJECT_NAME} -w /go/src/${PROJECT_NAME} ${BUILD_IMAGE} make docker-disk

release:portal build
	rm -rf release && mkdir -p release/stone/bin
	cp scripts/install.sh release/stone/
	cp -r scripts/systemd release/stone/
	cp stone release/stone/bin/
	cp tools/docker-disk release/stone/bin/
	cd release && tar zcvf stone-${MAJOR_VERSION}-${GIT_VERSION}.tar.gz stone && cd ..

clean:
	rm -rf stone

clean-docker-disk:
	rm -rf docker-disk

.PHONY: build
