// +build tools

package tools

// Manage tool dependencies via go.mod.
//
// https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
// https://github.com/golang/go/issues/25922
//
// nolint
import (
	_ "entgo.io/ent/cmd/ent"
	_ "github.com/abice/go-enum"
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "github.com/envoyproxy/protoc-gen-validate"
	_ "github.com/golang/mock/gomock"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/hexdigest/gowrap"
	_ "github.com/twitchtv/twirp"
	_ "golang.org/x/tools/cmd/goimports"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
