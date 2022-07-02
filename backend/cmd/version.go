package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

var revision = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print app version",
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {

		}
		fmt.Printf("File Browser %s\n", buildInfo.Main.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
