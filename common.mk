SHELL := /bin/bash
BASE_PATH := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION ?= $(shell git describe --tags --always --match=v* 2> /dev/null || \
           			cat $(CURDIR)/.version 2> /dev/null || echo v0)
VERSION_HASH = $(shell git rev-parse HEAD)

# printing
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

## help: Show this help
.PHONY: help
help:
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /' | sort


# tools
TOOLS_DIR := $(BASE_PATH)/tools
TOOLS_BIN := $(TOOLS_DIR)/bin
PATH := $(TOOLS_BIN):$(PATH)
export PATH

# go tools
.PHONY: go-deps
go-deps: $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/go.sum

go-enum=$(TOOLS_BIN)/go-enum
install-go-enum: $(TOOLS_BIN)/go-enum
$(TOOLS_BIN)/go-enum: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/go-enum: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/go-enum github.com/abice/go-enum

goimports=$(TOOLS_BIN)/goimports
install-goimports: $(TOOLS_BIN)/goimports
$(TOOLS_BIN)/goimports: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/goimports: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/goimports golang.org/x/tools/cmd/goimports

golangci-lint=$(TOOLS_BIN)/golangci-lint
install-golangci-lint: $(TOOLS_BIN)/golangci-lint
$(TOOLS_BIN)/golangci-lint: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/golangci-lint: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

mockgen=$(TOOLS_BIN)/mockgen
install-mockgen: $(TOOLS_BIN)/mockgen
$(TOOLS_BIN)/mockgen: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/mockgen: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/mockgen github.com/golang/mock/mockgen

buf=$(TOOLS_BIN)/buf
install-buf: $(TOOLS_BIN)/buf
$(TOOLS_BIN)/buf: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/buf: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/buf github.com/bufbuild/buf/cmd/buf

protoc=$(TOOLS_BIN)/protoc
install-protoc: $(TOOLS_BIN)/protoc
$(TOOLS_BIN)/protoc: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/protoc: tools/install_protoc.sh
	$Q cd ${TOOLS_DIR} && ./install_protoc.sh && touch $@

protoc-gen-validate=$(TOOLS_BIN)/protoc-gen-validate
install-protoc-gen-validate: $(TOOLS_BIN)/protoc-gen-validate
$(TOOLS_BIN)/protoc-gen-validate: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/protoc-gen-validate: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/protoc-gen-validate github.com/envoyproxy/protoc-gen-validate

protoc-gen-twirp=$(TOOLS_BIN)/protoc-gen-twirp
install-protoc-gen-twirp: $(TOOLS_BIN)/protoc-gen-twirp
$(TOOLS_BIN)/protoc-gen-twirp: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/protoc-gen-twirp: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/protoc-gen-twirp github.com/twitchtv/twirp/protoc-gen-twirp

gowrap=$(TOOLS_BIN)/gowrap
install-gowrap: $(TOOLS_BIN)/gowrap
$(TOOLS_BIN)/gowrap: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/gowrap: go-deps
	$Q cd ${TOOLS_DIR} && go build -o bin/gowrap github.com/hexdigest/gowrap/cmd/gowrap

# js tools
.PHONY: js-deps
js-deps: $(TOOLS_DIR)/node_modules/.modified
$(TOOLS_DIR)/node_modules/.modified: $(TOOLS_DIR)/package.json $(TOOLS_DIR)/yarn.lock
	$Q cd ${TOOLS_DIR} && yarn install
#	$Q find ${TOOLS_DIR}/node_modules -type f | xargs touch -am
	$Q touch $@

standard-version=$(TOOLS_BIN)/standard-version
install-standard-version: $(TOOLS_BIN)/standard-version
$(TOOLS_BIN)/standard-version: js-deps
	$Q ln -sf $(TOOLS_DIR)/node_modules/.bin/standard-version $@
	$Q touch $@

commitlint=$(TOOLS_BIN)/commitlint
install-commitlint: $(TOOLS_BIN)/commitlint
$(TOOLS_BIN)/commitlint: js-deps
	$Q ln -sf $(TOOLS_DIR)/node_modules/.bin/commitlint $@
	$Q touch $@
