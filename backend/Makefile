SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

include $(SELF_DIR)../common.mk

.DEFAULT_GOAL := test

MODULE = $(shell env GO111MODULE=on go list -m)
LDFLAGS += -X "main.revision=$(VERSION)"

## build: Build backend
.PHONY: build
build: | ; $(info $(M) building backend…)
	$Q go build  -ldflags '$(LDFLAGS)'

## test-backend: Run tests
.PHONY: test
test: | ; $(info $(M) running backend tests…)
	$Q go test -race -timeout 30s ./...

## generate: Run generators
.PHONY: generate
generate: install-go-enum install-mockgen | ; $(info $(M) running generators…)
	$Q go generate ./...
	$Q $(MAKE) fmt

## lint: Run all linters
.PHONY: lint
lint: install-golangci-lint | ; $(info $(M) running linters…)
	$Q $(golangci-lint) run

## fmt: Run formatting tools
.PHONY: fmt
fmt: install-goimports | ; $(info $(M) formatting source files…)
	$Q $(goimports) -local $(MODULE) -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

go.mod:
	go mod tidy
	go mod verify
go.sum: go.mod
