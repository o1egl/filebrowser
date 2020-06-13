.DEFAULT_GOAL := test
SHELL := /bin/bash
VERSION ?= $(shell git describe --tags --always --match=v* 2> /dev/null || \
           			cat $(CURDIR)/.version 2> /dev/null || echo v0)
VERSION_HASH = $(shell git rev-parse HEAD)

BASE_PATH := $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

TOOLS := $(BASE_PATH)/tools
TOOLS_BIN := $(TOOLS)/bin
PATH := $(TOOLS_BIN):$(PATH)
export PATH

MODULE = $(shell env GO111MODULE=on go list -m)
LDFLAGS += -X "main.revision=$(VERSION)"

# printing
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

## build-frontend: Build frontend
.PHONY: build-frontend
build-frontend: | ; $(info $(M) building frontend…)
	$Q cd frontend && npm install && npm run build

## build: Build backend
.PHONY: build-backend
build-backend: | ; $(info $(M) building backend…)
	$Q go build  -ldflags '$(LDFLAGS)'

## build: Build application
.PHONY: build
build: | build-frontend build-backend

## test: Run tests
.PHONY: test
test: | test-backend

## test-backend: Run backend tests
.PHONY: test-backend
test-backend: | ; $(info $(M) running backend tests…)
	$Q go test -race -timeout 30s ./...

## help: Show this help
.PHONY: help
help: Makefile
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /' | sort

## generate: Run generators
.PHONY: generate
generate: $(TOOLS_BIN)/go-enum $(TOOLS_BIN)/mockgen | ; $(info $(M) running generators…)
	$Q go generate ./...
	$Q $(MAKE) fmt

## lint: Run all linters
.PHONY: lint
lint: lint-backend lint-frontend lint-commits | ; $(info $(M) running linters…)

## lint-backend: Run backend linters
.PHONY: lint-backend
lint-backend: $(TOOLS_BIN)/golangci-lint | ; $(info $(M) running backend linters…)
	$Q $(TOOLS_BIN)/golangci-lint run

## lint-frontend: Run frontend linters
.PHONY: lint-frontend
lint-frontend: | ; $(info $(M) running frontend linters…)
	$Q cd frontend && npm install && npm run lint

## lint-commits: Run commit linters
.PHONY: lint-commits
lint-commits: $(TOOLS_BIN)/commitlint | ; $(info $(M) running commit linters…)
	$Q ./scripts/commitlint.sh

## fmt: Run formatting tools
.PHONY: fmt
fmt: $(TOOLS_BIN)/goimports | ; $(info $(M) formatting source files…)
	$Q $(TOOLS_BIN)/goimports -local $(MODULE) -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

## release: Bump applivation version
.PHONY: bump-version
bump-version: $(TOOLS_BIN)/standard-version | ; $(info $(M) bumping app version…)
	$Q ./scripts/bump_version.sh

## clean: Remove generated files
.PHONY: clean
clean:
	@# tools
	$Q rm -rf $(TOOLS_BIN)
	$Q rm -rf tools/node_modules
	@# frontend
	$Q rm -rf fronetend/node_modules
	$Q rm -rf fronetend/dist

# tools
$(TOOLS_BIN)/go-enum: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/go-enum: tools/go.mod tools/go.sum
	$Q cd tools && go build -o bin/go-enum github.com/abice/go-enum

$(TOOLS_BIN)/goimports: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/goimports: tools/go.mod tools/go.sum
	$Q cd tools && go build -o bin/goimports golang.org/x/tools/cmd/goimports

$(TOOLS_BIN)/golangci-lint: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/golangci-lint: tools/go.mod tools/go.sum
	$Q cd tools && go build -o bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint

$(TOOLS_BIN)/mockgen: export GOFLAGS = -mod=readonly
$(TOOLS_BIN)/mockgen: tools/go.mod tools/go.sum
	$Q cd tools && go build -o bin/mockgen github.com/golang/mock/mockgen

$(TOOLS_BIN)/npm_deps: tools/package.json tools/package-lock.json
	$Q cd tools && npm ci
	$Q find tools/node_modules/ -type f | xargs touch -am
	$Q mkdir -p $$(dirname $@) && touch $@

$(TOOLS_BIN)/standard-version: $(TOOLS_BIN)/npm_deps
	$Q ln -sf $(TOOLS)/node_modules/.bin/standard-version $@

$(TOOLS_BIN)/commitlint: $(TOOLS_BIN)/npm_deps
	$Q ln -sf $(TOOLS)/node_modules/.bin/commitlint $@

go.mod:
	go mod tidy
	go mod verify
go.sum: go.mod
