#Inspired from : https://github.com/littlemanco/boilr-makefile/blob/master/template/Makefile, https://github.com/geetarista/go-boilerplate/blob/master/Makefile, https://github.com/nascii/go-boilerplate/blob/master/GNUmakefile https://github.com/cloudflare/hellogopher/blob/master/Makefile
#PATH=$(PATH:):$(GOPATH)/bin
APP_NAME=docker-volume-gvfs
APP_VERSION=$(shell git describe --abbrev=0)
APP_USERREPO=github.com/sapk
APP_PACKAGE=$(APP_USERREPO)/docker-volume-gvfs

PLUGIN_NAME=sapk/$(APP_NAME)
PLUGIN_TAG=latest

GIT_HASH=$(shell git rev-parse --short HEAD)
GIT_BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
DATE := $(shell date -u '+%Y-%m-%d-%H%M-UTC')
PWD=$(shell pwd)

ARCHIVE=$(APP_NAME)-$(APP_VERSION)-$(GIT_HASH).tar.gz
#DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
LDFLAGS = \
  -s -w \
  -X main.Version=$(APP_VERSION}) -X main.Branch=$(GIT_BRANCH) -X main.Commit=$(GIT_HASH) -X main.BuildTime=$(DATE)

FAKE_GOPATH = $(PWD)/.gopath
FAKE_PACKAGE = $(FAKE_GOPATH)/src/$(APP_PACKAGE)

GO15VENDOREXPERIMENT=1
DOC_PORT = 6060
#GOOS=linux

ERROR_COLOR=\033[31;01m
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
WARN_COLOR=\033[33;01m


all: build compress done

build: deps format clean compile

push:  clean docker rootfs create enable
	@echo -e "$(OK_COLOR)==> push plugin ${PLUGIN_NAME}:${PLUGIN_TAG}$(NO_COLOR)"
	@docker plugin push ${PLUGIN_NAME}:${PLUGIN_TAG}
	
rootfs:
	@echo -e "$(OK_COLOR)==> create rootfs directory in ./plugin/rootfs$(NO_COLOR)"
	@mkdir -p ./plugin/rootfs
	@docker create --name tmp ${PLUGIN_NAME}:rootfs
	@docker export tmp | tar -x -C ./plugin/rootfs
	@echo -e "### copy config.json to ./plugin/$(NO_COLOR)"
	@cp config.json ./plugin/
	@docker rm -vf tmp
	
docker:
	@echo -e "$(OK_COLOR)==> Docker build image$(NO_COLOR)"
	@docker build -q -t ${PLUGIN_NAME}:rootfs -f .support/docker/Dockerfile /dev/null
	
create:
	@echo -e "$(OK_COLOR)==> Remove existing plugin ${PLUGIN_NAME}:${PLUGIN_TAG} if exists$(NO_COLOR)"
	@docker plugin rm -f ${PLUGIN_NAME}:${PLUGIN_TAG} || true
	@echo -e "$(OK_COLOR)==> Create new plugin ${PLUGIN_NAME}:${PLUGIN_TAG} from ./plugin$(NO_COLOR)"
	@docker plugin create ${PLUGIN_NAME}:${PLUGIN_TAG} ./plugin

enable:
	@echo -e "$(OK_COLOR)==> Enable plugin ${PLUGIN_NAME}:${PLUGIN_TAG}$(NO_COLOR)"
	@docker plugin enable ${PLUGIN_NAME}:${PLUGIN_TAG}

set-build:
	@if [ ! -d $(PWD)/.gopath/src/$(APP_USERREPO) ]; then mkdir -p $(PWD)/.gopath/src/$(APP_USERREPO); fi
	@if [ ! -d $(PWD)/.gopath/src/$(APP_PACKAGE) ]; then ln -s $(PWD) $(PWD)/.gopath/src/$(APP_PACKAGE); fi

compile: set-build
	@echo -e "$(OK_COLOR)==> Building...$(NO_COLOR)"
	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) go build -v -ldflags "$(LDFLAGS)"

release: clean set-build deps format
	@mkdir build
	@echo -e "$(OK_COLOR)==> Building for linux 32 ...$(NO_COLOR)"
	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -o build/${APP_NAME}-linux-386 -ldflags "$(LDFLAGS)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-386 || upx-ucl --brute  build/${APP_NAME}-linux-386 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for linux 64 ...$(NO_COLOR)"
	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) GO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/${APP_NAME}-linux-amd64 -ldflags "$(LDFLAGS)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-amd64 || upx-ucl --brute  build/${APP_NAME}-linux-amd64 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for linux arm ...$(NO_COLOR)"
	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -o build/${APP_NAME}-linux-armv6 -ldflags "$(LDFLAGS)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-linux-armv6 || upx-ucl --brute  build/${APP_NAME}-linux-armv6 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for darwin32 ...$(NO_COLOR)"
	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) CGO_ENABLED=0 GOOS=darwin GOARCH=386 go build -o build/${APP_NAME}-darwin-386 -ldflags "$(LDFLAGS)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-darwin-386 || upx-ucl --brute  build/${APP_NAME}-darwin-386 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Building for darwin64 ...$(NO_COLOR)"
	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/${APP_NAME}-darwin-amd64 -ldflags "$(LDFLAGS)"
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute  build/${APP_NAME}-darwin-amd64 || upx-ucl --brute  build/${APP_NAME}-darwin-amd64 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

#	@echo -e "$(OK_COLOR)==> Building for win32 ...$(NO_COLOR)"
#	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o build/${APP_NAME}-win-386 -ldflags "$(LDFLAGS)"
#	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
#	@upx --brute  build/${APP_NAME}-win-386 || upx-ucl --brute  build/${APP_NAME}-win-386 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

#	@echo -e "$(OK_COLOR)==> Building for win64 ...$(NO_COLOR)"
#	cd $(FAKE_PACKAGE) && GOPATH=$(FAKE_GOPATH) CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/${APP_NAME}-win-amd64 -ldflags "$(LDFLAGS)"
#	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
#	@upx --brute  build/${APP_NAME}-win-amd64 || upx-ucl --brute  build/${APP_NAME}-win-amd64 || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

	@echo -e "$(OK_COLOR)==> Archiving ...$(NO_COLOR)"
	@tar -zcvf build/$(ARCHIVE) LICENSE README.md build/

clean:
	@if [ -x docker-volume-gvfs ]; then rm docker-volume-gvfs; fi
	@if [ -d build ]; then rm -R build; fi
	@if [ -d $(FAKE_GOPATH) ]; then rm -R $(FAKE_GOPATH); fi
	@rm -rf ./plugin

compress:
	@echo -e "$(OK_COLOR)==> Trying to compress binary ...$(NO_COLOR)"
	@upx --brute docker-volume-gvfs || upx-ucl --brute docker-volume-gvfs || echo -e "$(WARN_COLOR)==> No tools found to compress binary.$(NO_COLOR)"

format:
	@echo -e "$(OK_COLOR)==> Formatting...$(NO_COLOR)"
	go fmt . ./gvfs/...
#go fmt ./...

test: deps format
	@echo -e "$(OK_COLOR)==> Running tests...$(NO_COLOR)"
	go vet . ./gvfs/... || true
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./gvfs/drivers
	go tool cover -html=coverage.out -o coverage.html

docs:
	@echo -e "$(OK_COLOR)==> Serving docs at http://localhost:$(DOC_PORT).$(NO_COLOR)"
	@godoc -http=:$(DOC_PORT)

lint: dev-deps
	gometalinter --deadline=5m --concurrency=2 --vendor ./...

dev-deps:
	@echo -e "$(OK_COLOR)==> Installing developement dependencies...$(NO_COLOR)"
	@go get github.com/nsf/gocode
	@go get github.com/alecthomas/gometalinter
	@go get github.com/dpw/vendetta #Vendoring
	@$(GOPATH)/bin/gometalinter --install > /dev/null

update-dev-deps:
	@echo -e "$(OK_COLOR)==> Installing/Updating developement dependencies...$(NO_COLOR)"
	go get -u github.com/nsf/gocode
	go get -u github.com/alecthomas/gometalinter
	go get -u github.com/dpw/vendetta #Vendoring
	$(GOPATH)/bin/gometalinter --install --update

deps:
	@echo -e "$(OK_COLOR)==> Installing dependencies ...$(NO_COLOR)"
	@git submodule update --init --recursive
# @$(GOPATH)/bin/vendetta -n $(APP_PACKAGE)
#	@go get -d -v ./...

update-deps: dev-deps
	@echo -e "$(OK_COLOR)==> Updating all dependencies ...$(NO_COLOR)"
	$(GOPATH)/bin/vendetta -n $(APP_PACKAGE) -u
#@go get -d -v -u ./...


done:
	@echo -e "$(OK_COLOR)==> Done.$(NO_COLOR)"

.PHONY: all build compile clean compress format test docs lint dev-deps update-dev-deps deps update-deps done
