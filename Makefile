VERSION=0.0.1
COMMIT=$(shell git log -q -1 | head -n 1 | cut -f2 -d' ')
TAG=none

GO15VENDOREXPERIMENT=1

all: deps compile

compile:
	go build -ldflags "-X main.Version=v${VERSION} -X main.Tag=${TAG} -X main.Commit=${COMMIT}"
	#go build

format:
	gofmt -s -w -l .

deps:
	go get -v
#	go get -v ./...
