VERSION=0.0.1
COMMIT=$(shell git log -q -1 | head -n 1 | cut -f2 -d' ')
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

GO15VENDOREXPERIMENT=1
GOPATH ?= $(GOPATH:):./vendor

all: deps compile

compile:
	go build -ldflags "-X main.version=v${VERSION} -X main.branch=${BRANCH} -X main.commit=${COMMIT}"
	#go build

format:
	gofmt -s -w -l .

deps:
	go get -dv
#	go get -v ./...
