package server

import (
	"context" //nolint:gosec
	"os"
	"os/signal"
	"syscall"

	"github.com/filebrowser/filebrowser/v3/cmd"
	"github.com/filebrowser/filebrowser/v3/config"
	"github.com/filebrowser/filebrowser/v3/domain"
	"github.com/filebrowser/filebrowser/v3/log"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Command with command line flags and env
type Command struct {
	Config string `long:"config" short:"c" env:"CONFIG" description:"path to config.yaml"`

	cmd.CommonOpts
}

// Execute runs file browser server
func (s *Command) Execute(_ []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warnf("interrupt signal")
		cancel()
	}()

	cfg := config.Default()
	if s.Config != "" {
		var err error
		cfg, err = config.Load(ctx, config.FileLoader(s.Config))
		if err != nil {
			log.Fatalf("failed to load config, %+v", err)
			return err
		}
	}

	app, err := InitializeServer(ctx, cfg, domain.Version(s.Revision))
	if err != nil {
		log.Fatalf("failed to setup application, %+v", err)
		return err
	}

	if err = app.run(ctx); err != nil {
		log.Fatalf("terminated with error %+v", err)
		return err
	}
	log.Infof("terminated")
	return nil
}
