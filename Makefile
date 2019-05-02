GOOS ?= $(shell go env GOOS)
DEFAULT_ARCH := amd64
BUILD_TIME = $(shell date -u +%Y%m%d.%H%M%S)
VERSION = $(shell cat version)
LDFLAGS = -ldflags "-s -w -extldflags '-static' -X main.build=$(BUILD_TIME) -X main.version=$(VERSION)"
BINARY = goapp-$(VERSION)-$(GOOS)-amd64
OUTPUT ?= build/${BINARY}$(ext) 
$(eval ext := $(if $(filter $(GOOS),windows),.exe))

default: build

.PHONY: build
build: 
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${DEFAULT_ARCH} go build -a ${LDFLAGS} -o ${OUTPUT} .

clean:
	rm -rf build/*
