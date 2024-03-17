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
	_ "github.com/abice/go-enum"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/swaggo/swag/cmd/swag"
	_ "go.uber.org/mock/mockgen"
	_ "mvdan.cc/gofumpt"
)
