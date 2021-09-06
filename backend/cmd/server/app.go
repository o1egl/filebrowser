package server

import (
	"context"

	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/rest/api"
)

// serverApp holds all active objects
type serverApp struct {
	*ServerCommand
	restSrv    *api.Server
	terminated chan struct{}
}

func NewServerApp(srvCmd *ServerCommand, restSrv *api.Server) (*serverApp, error) {
	return &serverApp{
		ServerCommand: srvCmd,
		restSrv:       restSrv,
		terminated:    make(chan struct{}),
	}, nil
}

// Run all application objects
func (a *serverApp) run(ctx context.Context) error {
	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		log.Warnf("shutdown initiated")
		a.restSrv.Shutdown()
	}()

	a.restSrv.Run()

	close(a.terminated)
	return nil
}

// Wait for application completion (termination)
func (a *serverApp) Wait() {
	<-a.terminated
}
