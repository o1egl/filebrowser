package api

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/filebrowser/filebrowser/config"
	_ "github.com/filebrowser/filebrowser/docs"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.Config) *Server {
	r := chi.NewRouter()

	fileService := NewFileService()

	serverURI, err := url.Parse(cfg.PublicAddress())
	if err != nil {
		slog.Error("Failed to parse server uri: %w", err)
		os.Exit(1)
	}

	r.Route(cfg.BasePath, func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Route("/local", func(r chi.Router) {
				r.Post("/login", nil)
			})
			r.Delete("/logout", nil)
		})
		r.Route("/api", func(r chi.Router) {
			r.Route("/v1", func(r chi.Router) {
				r.Route("/files", func(r chi.Router) {
					r.Get("/list", fileService.List)
					r.Patch("/rename", fileService.Rename)
					r.Patch("/move", fileService.Move)
					r.Post("/copy", fileService.Copy)
					r.Delete("/delete", fileService.Delete)
				})
			})
			r.Handle("/swagger/*", swaggerHandler(serverURI))
		})
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
