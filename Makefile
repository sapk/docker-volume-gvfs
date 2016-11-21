VERSION = 0.0.1
GO15VENDOREXPERIMENT=1

all: deps compile

compile: goxc

format:
	gofmt -s -w -l .

deps:
	go get -v
