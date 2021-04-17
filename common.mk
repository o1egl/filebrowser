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
TOOLS_GO_DEPS := $(TOOLS_DIR)/go.mod $(TOOLS_DIR)/go.sum

go-enum=$(TOOLS_BIN)/go-enum
$(go-enum): $(TOOLS_GO_DEPS)
	$Q cd ${TOOLS_DIR} && go build -o $@ github.com/abice/go-enum

goimports=$(TOOLS_BIN)/goimports
$(goimports): $(TOOLS_GO_DEPS)
	$Q cd ${TOOLS_DIR} && go build -o $@ golang.org/x/tools/cmd/goimports

golangci-lint=$(TOOLS_BIN)/golangci-lint
$(golangci-lint): $(TOOLS_GO_DEPS)
	$Q cd ${TOOLS_DIR} && go build -o $@ github.com/golangci/golangci-lint/cmd/golangci-lint

mockgen=$(TOOLS_BIN)/mockgen
$(mockgen): $(TOOLS_GO_DEPS)
	$Q cd ${TOOLS_DIR} && go build -o $@ github.com/golang/mock/mockgen

gowrap=$(TOOLS_BIN)/gowrap
$(gowrap): $(TOOLS_GO_DEPS)
	$Q cd ${TOOLS_DIR} && go build -o $@ github.com/hexdigest/gowrap/cmd/gowrap

ent=$(TOOLS_BIN)/ent
$(ent): $(TOOLS_GO_DEPS)
	$Q cd ${TOOLS_DIR} && go build -o $@ entgo.io/ent/cmd/ent

# js tools
TOOLS_JS_DEPS: $(TOOLS_DIR)/node_modules/.modified
$(TOOLS_JS_DEPS): $(TOOLS_DIR)/package.json $(TOOLS_DIR)/yarn.lock
	$Q cd ${TOOLS_DIR} && yarn install
#	$Q find ${TOOLS_DIR}/node_modules -type f | xargs touch -am
	$Q touch -am $@

standard-version=$(TOOLS_BIN)/standard-version
$(standard-version): $(TOOLS_JS_DEPS)
	$Q ln -sf $(TOOLS_DIR)/node_modules/.bin/standard-version $@
	$Q touch -am $@

commitlint=$(TOOLS_BIN)/commitlint
$(commitlint): $(TOOLS_JS_DEPS)
	$Q ln -sf $(TOOLS_DIR)/node_modules/.bin/commitlint $@
	$Q touch -am $@
