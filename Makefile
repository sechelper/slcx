EXECUTABLES = git go pwd mkdir
K := $(foreach exec,$(EXECUTABLES),\
        $(if $(shell which $(exec)),some string,$(error "No $(exec) in PATH)))

ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))/build

BINARY=slcx
VERSION=v1.1
BUILD=`git rev-parse --short HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=386 amd64
BUILD_DIR=target

# Setup linker flags option for build that interoperate with variable names in src code
LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD}"

default: build

all: clean build_all

build:
	go build ${LDFLAGS} -o ${BINARY}

build_all:
	mkdir -p ${BUILD_DIR}
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build -v -o ${BUILD_DIR}/$(BINARY)_${VERSION}_$(GOOS)_$(GOARCH))))


# Remove only what we've created
clean:
	rm -rf ${BUILD_DIR}

.PHONY: check clean install build_all all
