SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

include $(SELF_DIR)../common.mk

.DEFAULT_GOAL := test

MODULE = $(shell env GO111MODULE=on go list -m)
LDFLAGS += -X "main.revision=$(VERSION)"

## build: Build backend
.PHONY: build
build: | ; $(info $(M) building…)
	$Q go build  -ldflags '$(LDFLAGS)'

## test: Run tests
.PHONY: test
test: | ; $(info $(M) running tests…)
	$Q go test -race -timeout 30s ./...

## generate: Run generators
.PHONY: generate
generate: $(go-enum) $(mockgen) $(ent) | ; $(info $(M) running generators…)
	$Q go generate ./...
	$Q $(MAKE) fmt

## lint: Run all linters
.PHONY: lint
lint: $(golangci-lint) | ; $(info $(M) running linters…)
	$Q $(golangci-lint) run

## fmt: Run formatting tools
.PHONY: fmt
fmt: $(goimports) | ; $(info $(M) formatting source files…)
	$Q $(goimports) -local $(MODULE) -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

go.mod:
	go mod tidy
	go mod verify
go.sum: go.mod
