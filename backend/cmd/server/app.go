package server

import (
	"context"

	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/rest/api"
	httpServer "github.com/filebrowser/filebrowser/v3/server"
)

// app holds all active objects
type app struct {
	restSrv    *api.Server
	srv        *httpServer.Server
	terminated chan struct{}
}

func newApp(restSrv *api.Server, srv *httpServer.Server) (*app, error) {
	return &app{
		restSrv:    restSrv,
		srv:        srv,
		terminated: make(chan struct{}),
	}, nil
}

// Run all application objects
func (a *app) run(ctx context.Context) error {
	go func() {
		// shutdown on context cancellation
		<-ctx.Done()
		log.Warnf("shutdown initiated")
		//a.restSrv.Shutdown()
		a.srv.Shutdown()
	}()

	//a.restSrv.Run()
	a.srv.Run()

	close(a.terminated)
	return nil
}

// Wait for application completion (termination)
func (a *app) Wait() {
	<-a.terminated
}
