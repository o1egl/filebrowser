package cmd

import (
	"context"
	"github.com/bep/simplecobra"
	"github.com/filebrowser/filebrowser/api"
	"github.com/filebrowser/filebrowser/config"
	"github.com/filebrowser/filebrowser/logger"
	"golang.org/x/exp/slog"
	"io"
	"os"
)

type serveCommand struct {
	cfg        *config.Config
	configPath string
}

func newServeCommand() *serveCommand {
	return &serveCommand{}
}

func (s *serveCommand) Name() string {
	return "serve"
}

func (s *serveCommand) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	cmd.Short = "Start the filebrowser server"
	cmd.Flags().StringVarP(&s.configPath, "config", "c", "config.yaml", "config file path")
	return nil
}

func (s *serveCommand) PreRun(cd, runner *simplecobra.Commandeer) (err error) {
	s.cfg, err = config.FromFile(s.configPath)
	if err != nil {
		return err
	}
	if err = initLogger(s.cfg.Log); err != nil {
		return err
	}
	return nil
}

func (s *serveCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	server := api.NewServer(s.cfg)

	return server.Run(ctx)
}

func (s *serveCommand) Commands() []simplecobra.Commander {
	return nil
}

func initLogger(cfg config.Log) error {
	var lvl slog.Level
	switch cfg.Level {
	case config.LogLevelDebug:
		lvl = slog.LevelDebug
	case config.LogLevelInfo:
		lvl = slog.LevelInfo
	case config.LogLevelWarn:
		lvl = slog.LevelWarn
	case config.LogLevelError:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	var logWriter io.Writer
	switch cfg.Output {
	case config.LogOutputStdout:
		logWriter = os.Stdout
	case config.LogOutputStderr:
		logWriter = os.Stderr
	default:
		logWriter = os.Stdout
	}

	handlerOpts := &slog.HandlerOptions{
		Level: lvl,
	}

	var handler slog.Handler
	switch cfg.Format {
	case config.LogFormatText:
		handler = slog.NewTextHandler(logWriter, handlerOpts)
	case config.LogFormatJson:
		handler = slog.NewJSONHandler(logWriter, handlerOpts)
	default:
		handler = slog.NewTextHandler(logWriter, handlerOpts)
	}

	log := slog.New(logger.NewContextHandler(handler))
	slog.SetDefault(log)

	return nil
}
