package main

import (
	"fmt"
	"os"

	"github.com/jessevdk/go-flags"

	"github.com/filebrowser/filebrowser/v3/cmd"
	"github.com/filebrowser/filebrowser/v3/log"
)

// Opts with all cli commands and flags
type Opts struct {
	ServerCmd cmd.ServerCommand `command:"server"`

	Log LogGroup `group:"log" namespace:"log" env-namespace:"LOG"`

	ServerURL    string `long:"url" env:"SERVER_URL" required:"true" description:"url to file browser"`
	SharedSecret string `long:"secret" env:"SECRET" required:"true" description:"shared secret key"`
	Version      func() `long:"version" description:"print version"`
}

// LogGroup defines options group log params
type LogGroup struct {
	Level  string `long:"level" env:"LEVEL" default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"critical" choice:"fatal" description:"log level"` //nolint:lll,staticcheck
	Format string `long:"format" env:"FORMAT" default:"plain" choice:"plain" choice:"json" description:"log format"`                                                           //nolint:lll,staticcheck
}

var revision = "unknown"

func main() {
	var opts Opts
	opts.Version = func() {
		fmt.Printf("File Browser %s\n", revision)
		os.Exit(0)
	}
	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		if err := setupLogger(opts.Log.Level, opts.Log.Format); err != nil {
			fmt.Printf("[ERROR] failed to initialize logger: %s", err)
			os.Exit(1)
		}
		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
		c := command.(cmd.CommonOptionsCommander)
		c.SetCommon(cmd.CommonOpts{
			ServerURL:    opts.ServerURL,
			SharedSecret: opts.SharedSecret,
			Revision:     revision,
		})
		for _, entry := range c.HandleDeprecatedFlags() {
			log.Warnf("--%s is deprecated and will be removed in v%s, please use --%s instead",
				entry.Old, entry.RemoveVersion, entry.New)
		}
		err := c.Execute(args)
		if err != nil {
			log.Errorf("failed with %+v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			p.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func setupLogger(logLevel, logFormat string) error {
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}
	format, err := log.ParseFormat(logFormat)
	if err != nil {
		return err
	}
	logger, err := log.NewLogger(log.Configuration{
		LogLevel: level,
		Format:   format,
		Output:   os.Stdout,
	})
	if err != nil {
		return err
	}
	log.DefaultLogger = logger
	return nil
}
