//go:build tools
// +build tools

package tools

// Manage tool dependencies via go.mod.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
//
// nolint
import (
	_ "connectrpc.com/connect/cmd/protoc-gen-connect-go"
	_ "github.com/abice/go-enum"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "go.uber.org/mock/mockgen"
	_ "golang.org/x/tools/cmd/goimports"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
