package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/filebrowser/filebrowser/config"
	"github.com/filebrowser/filebrowser/domain"
	"github.com/filebrowser/filebrowser/hash"
	"github.com/filebrowser/filebrowser/log"
	"github.com/filebrowser/filebrowser/storage"
)

type Service interface {
	Start() error
	Stop(ctx context.Context) error
}

type App struct {
	cfg      *config.Config
	logger   log.Logger
	store    storage.Store
	hasher   hash.Hasher
	services []Service
}

func NewApp(
	cfg *config.Config,
	logger log.Logger,
	store storage.Store,
	hasher hash.Hasher,
	services []Service,
) *App {
	return &App{
		cfg:      cfg,
		logger:   logger,
		store:    store,
		services: services,
		hasher:   hasher,
	}
}

// Start all application services
func (a *App) Start() error {
	if err := a.createSystemUsers(); err != nil {
		return err
	}

	return a.startServices()
}

func (a *App) createSystemUsers() error {
	pwd, err := a.hasher.Password("admin")
	if err != nil {
		return err
	}

	users := []*domain.User{
		{
			Username: "admin",
			Password: pwd,
			Home: domain.HomeVolume{
				Path: "/",
				Permissions: domain.VolumePermissions{
					Read:   true,
					Create: true,
					Modify: true,
					Delete: true,
					Share:  true,
				},
			},
			Name:    "Admin",
			Locale:  "en",
			IsAdmin: true,
		},
		{
			Username: "guest",
			Home: domain.HomeVolume{
				Path: path.Join("/", a.cfg.Auth.Anonymous.Home.Path),
				Permissions: domain.VolumePermissions{
					Read:   a.cfg.Auth.Anonymous.Home.Permissions.Read,
					Create: a.cfg.Auth.Anonymous.Home.Permissions.Create,
					Modify: a.cfg.Auth.Anonymous.Home.Permissions.Modify,
					Delete: a.cfg.Auth.Anonymous.Home.Permissions.Delete,
					Share:  a.cfg.Auth.Anonymous.Home.Permissions.Share,
				},
			},
			Name:         "Guest",
			Locale:       "en",
			LockPassword: true,
		},
	}

	for _, user := range users {
		if err := a.createUserIfNoExist(user); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) createUserIfNoExist(user *domain.User) error {
	ctx := context.Background()
	existingUser, err := a.store.GetUserByUsername(ctx, user.Username)
	if err != nil && !domain.IsNotFoundError(err) {
		return err
	}

	// user already exist
	if existingUser != nil {
		return nil
	}

	a.logger.Debug(fmt.Sprintf("Creating %s user", user.Username))
	return a.store.CreateUser(ctx, user)
}

func (a *App) startServices() error {
	wait := make(chan struct{})

	// run services
	errs := make(chan error, len(a.services))
	for _, svc := range a.services {
		svc := svc
		go func() {
			if err := svc.Start(); err != nil {
				errs <- err
			}
		}()
	}

	// catch interrupt signal
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		a.logger.Warn("interrupt signal received")
		close(wait)
	}()

	select {
	case err := <-errs:
		return err
	case <-wait:
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelFn()
		return a.Stop(ctx)
	}
}

// Stop application
func (a *App) Stop(ctx context.Context) error {
	wg, wgCtx := errgroup.WithContext(ctx)
	for _, svc := range a.services {
		svc := svc
		wg.Go(func() error {
			return svc.Stop(wgCtx)
		})
	}

	done := make(chan struct{})
	errCh := make(chan error)
	go func() {
		if err := wg.Wait(); err != nil {
			errCh <- err
			return
		}
		close(done)
	}()

	select {
	case <-done:
		return nil
	case err := <-errCh:
		return err
	case <-ctx.Done():
		return fmt.Errorf("failed to gracefuly shutdown: %w", ctx.Err())
	}
}
