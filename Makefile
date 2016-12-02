VERSION=0.0.1
COMMIT=$(shell git log -q -1 | head -n 1 | cut -f2 -d' ')
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

GO15VENDOREXPERIMENT=1
GOPATH ?= $(GOPATH:):./vendor
#GOOS=linux


all: deps compile

compile:
	#go build -ldflags "-X main.version=v${VERSION} -X main.branch=${BRANCH} -X main.commit=${COMMIT}"
	go build -ldflags "-s -w -X main.version=v${VERSION} -X main.branch=${BRANCH} -X main.commit=${COMMIT}"
	#go build

compress:
	upx --brute docker-volume-gvfs || upx-ucl --brute docker-volume-gvfs

format:
	gofmt -s -w -l .

deps:
	go get -d -v
#	go get -v ./...
