package main

import (
	"github.com/filebrowser/filebrowser/cmd"
	"golang.org/x/exp/slog"
	"os"
)

func main() {
	err := cmd.Execute(os.Args[1:])
	if err != nil {
		slog.Error(err.Error())
	}
}
