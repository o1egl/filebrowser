package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/filebrowser/filebrowser/config"
	"github.com/filebrowser/filebrowser/gen/proto/filebrowser/v1/filebrowserv1connect"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	r := chi.NewRouter()

	fileService := NewFileService()
	userService := NewUserService()

	r.Route("/api", func(r chi.Router) {
		fileServicePath, fileServiceHandler := filebrowserv1connect.NewFileServiceHandler(fileService)
		r.Handle(fileServicePath, fileServiceHandler)

		userServicePath, userServiceHandler := filebrowserv1connect.NewUserServiceHandler(userService)
		r.Handle(userServicePath, userServiceHandler)
	})

	return &Server{
		httpServer: &http.Server{
			Addr:    net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
			Handler: r,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	listenErrCh := make(chan error)
	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			listenErrCh <- err
		}
	}()

	shutdownErrCh := make(chan error)
	go func() {
		<-ctx.Done()
		shutdownCtx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		shutdownErrCh <- s.httpServer.Shutdown(shutdownCtx)
	}()

	select {
	case err := <-listenErrCh:
		return fmt.Errorf("failed to start http server: %w", err)
	case err := <-shutdownErrCh:
		if err != nil {
			return fmt.Errorf("failed to shutdown http server: %w", err)
		}
		slog.Info("HTTP server shutdown successfully")
		return nil
	}
}
