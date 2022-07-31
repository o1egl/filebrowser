package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/filebrowser/filebrowser/cmd/server"
	"github.com/filebrowser/filebrowser/config"
	"github.com/filebrowser/filebrowser/hash"
	"github.com/filebrowser/filebrowser/storage"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

// versionCmd represents the version command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start filebrowser server",
	Run: func(command *cobra.Command, args []string) {
		ctx, cancelFn := context.WithCancel(context.Background())
		defer cancelFn()
		_ = ctx

		logger, err := NewLogger()
		if err != nil {
			panic(err)
		}
		cfg, err := config.LoadWithDefaults()
		if err != nil {
			logger.Fatal("Failed to load config", zap.Error(err))
		}

		hasher := hash.NewHasher(cfg.Secret)
		store, err := storage.New(cfg.Store)
		if err != nil {
			logger.Fatal("Failed to initialize storage", zap.Error(err))
		}

		app := server.NewApp(cfg, logger, store, hasher, nil)
		if err := app.Start(); err != nil {
			logger.Fatal("Failed to start the server", zap.Error(err))
		}
	},
}
