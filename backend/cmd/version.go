package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/filebrowser/filebrowser/buildinfo"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print app version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", buildinfo.Revision())
	},
}
