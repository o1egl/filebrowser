package main

import (
	"os"

	"github.com/filebrowser/filebrowser/cmd"
	"golang.org/x/exp/slog"
)

func main() {
	err := cmd.Execute(os.Args[1:])
	if err != nil {
		slog.Error(err.Error())
	}
}
