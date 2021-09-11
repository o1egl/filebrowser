package main

import (
	"fmt"
	"os"

	"github.com/filebrowser/filebrowser/v3/cmd/server"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"

	"github.com/filebrowser/filebrowser/v3/cmd"
	"github.com/filebrowser/filebrowser/v3/log"
)

// Opts with all cli commands and flags
type Opts struct {
	ServerCmd   server.Command      `command:"server"`
	PasswordCmd cmd.PasswordCommand `command:"password"`

	Log     LogGroup `group:"log" namespace:"log" env-namespace:"LOG"`
	Version func()   `long:"version" description:"print version"`
}

// LogGroup defines options group log params
type LogGroup struct {
	Level  string `long:"level" env:"LEVEL" default:"info" choice:"debug" choice:"info" choice:"warn" choice:"error" choice:"critical" choice:"fatal" description:"log level"` //nolint:lll,staticcheck
	Format string `long:"format" env:"FORMAT" default:"plain" choice:"plain" choice:"json" description:"log format"`                                                           //nolint:lll,staticcheck
	Out    string `long:"out" env:"OUT" choice:"stdout" choice:"file" default:"stdout" description:"output"`
	File   struct {
		Name       string `long:"name" env:"NAME" default:"./var/log/access.log" description:"path to access.log file"`
		MaxSize    int    `long:"size" env:"SIZE" default:"500" description:"maximum size in megabytes"`
		MaxBackups int    `long:"backups" env:"BACKUPS" default:"3" description:"maximum number of old log files"`
		MaxAge     int    `long:"age" env:"AGE" default:"30" description:"maximum number of days to retain"`
		Compress   bool   `long:"compress" env:"COMPRESS" description:"compress rotated log files"`
	} `group:"file" namespace:"file" env-namespace:"FILE"`
}

var revision = "unknown"

func main() {
	var opts Opts
	p := flags.NewParser(&opts, flags.Default)

	opts.Version = func() {
		fmt.Printf("File Browser %s\n", revision)
		os.Exit(0)
	}

	p.CommandHandler = func(command flags.Commander, args []string) error {
		if err := setupLogger(opts.Log); err != nil {
			fmt.Printf("[ERROR] failed to initialize logger: %s", err)
			os.Exit(1)
		}
		// commands implements CommonOptionsCommander to allow passing set of extra options defined for all commands
		c := command.(cmd.CommonOptionsCommander)
		c.SetCommon(cmd.CommonOpts{
			Revision: revision,
		})
		for _, entry := range c.HandleDeprecatedFlags() {
			log.Warnf("--%s is deprecated and will be removed in v%s, please use --%s instead",
				entry.Old, entry.RemoveVersion, entry.New)
		}
		err := c.Execute(args)
		if err != nil {
			log.Errorf("failed with %v", err)
		}
		return err
	}

	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println(err)
			p.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}
}

func setupLogger(group LogGroup) error {
	level, err := log.ParseLevel(group.Level)
	if err != nil {
		return err
	}
	format, err := log.ParseFormat(group.Format)
	if err != nil {
		return err
	}
	out, err := loggerOut(group)
	if err != nil {
		return err
	}
	logger, err := log.NewLogger(log.Configuration{
		LogLevel: level,
		Format:   format,
		Output:   out,
	})
	if err != nil {
		return err
	}
	log.DefaultLogger = logger
	return nil
}

func loggerOut(group LogGroup) (log.WriteSyncer, error) {
	switch group.Out {
	case "stdout":
		return os.Stdout, nil
	case "file":
		fileWriter := log.NewFileWriter(log.FileWriterConfig{
			Filename:   group.File.Name,
			MaxSize:    group.File.MaxSize,
			MaxAge:     group.File.MaxAge,
			MaxBackups: group.File.MaxBackups,
			Compress:   group.File.Compress,
		})
		return fileWriter, nil
	default:
		return nil, errors.Errorf("unsupported log out %s", group.Out)
	}
}
