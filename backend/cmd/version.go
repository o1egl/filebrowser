package cmd

import (
	"context"
	"fmt"
	"github.com/bep/simplecobra"
)

var version = "dev"

func newVersionCmd() *simpleCommand {
	return &simpleCommand{
		name:  "version",
		short: "Print the version number",
		run: func(ctx context.Context, cd *simplecobra.Commandeer, rootCmd *rootCommand, args []string) error {
			fmt.Println(version)
			return nil
		},
	}
}
