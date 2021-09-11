package api

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/didip/tollbooth/v6"
	"github.com/filebrowser/filebrowser/v3/config"
	"github.com/filebrowser/filebrowser/v3/domain"
	"github.com/filebrowser/filebrowser/v3/hash"
	"github.com/filebrowser/filebrowser/v3/service/filebrowser"
	"github.com/filebrowser/filebrowser/v3/store"
	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pkgz/auth"

	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/rest/middleware"
)

const (
	ShutdownTimeout      = 10 * time.Second
	CuncurrentRequests   = 1000
	AdminCurrentRequests = 10

	AuthRouterLimiter = 5

	PublicRoutesTimeout = 5 * time.Second
	StaticRouterLimiter = 100

	ProtectedRoutesTimeout = 30 * time.Second
	ProtectedRouterLimiter = 10

	AdminRoutesTimeout = 30 * time.Second
	AdminRouterLimiter = 10
)

type Server struct {
	cfg            *config.Config
	fileBrowserSvc filebrowser.Service
	authenticator  *auth.Service
	userStore      store.UserStore
	hasher         hash.Hasher
	version        domain.Version

	httpsServer *http.Server
	httpServer  *http.Server
	lock        sync.Mutex
}

func NewServer(
	cfg *config.Config,
	fileBrowserSvc filebrowser.Service,
	authenticator *auth.Service,
	hasher hash.Hasher,
	userStore store.UserStore,
	version domain.Version,
) *Server {
	return &Server{
		cfg:            cfg,
		fileBrowserSvc: fileBrowserSvc,
		authenticator:  authenticator,
		userStore:      userStore,
		hasher:         hasher,
		version:        version,
	}
}

func (s *Server) Run() {
	httpAddr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	httpsAddr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.SSL.Port)
	switch s.cfg.Server.SSL.Mode {
	case config.SSLModeNone:
		log.Infof("activate http rest server on %s", httpAddr)

		s.lock.Lock()
		s.httpServer = s.makeHTTPServer(httpAddr, s.newEngine())
		s.httpServer.ErrorLog = log.ToStdLogger(log.DefaultLogger, log.LevelWarn)
		s.lock.Unlock()

		err := s.httpServer.ListenAndServe()
		log.Warnf("http server terminated, %s", err)
	case config.SSLModeStatic:
		log.Infof("activate https server in 'static' mode on %s", httpsAddr)

		s.lock.Lock()
		s.httpsServer = s.makeHTTPSServer(httpsAddr, s.newEngine())
		s.httpsServer.ErrorLog = log.ToStdLogger(log.DefaultLogger, log.LevelWarn)

		s.httpServer = s.makeHTTPServer(httpAddr, s.httpToHTTPSRouter())
		s.httpServer.ErrorLog = log.ToStdLogger(log.DefaultLogger, log.LevelWarn)
		s.lock.Unlock()

		go func() {
			log.Infof("activate http redirect server on %s", httpAddr)
			err := s.httpServer.ListenAndServe()
			log.Warnf("http redirect server terminated, %s", err)
		}()

		err := s.httpsServer.ListenAndServeTLS(s.cfg.Server.SSL.Cert, s.cfg.Server.SSL.Key)
		log.Warnf("https server terminated, %s", err)
	case config.SSLModeAuto:
		log.Infof("activate https server in 'auto' mode on %s", httpsAddr)

		m := s.makeAutocertManager()
		s.lock.Lock()
		s.httpsServer = s.makeHTTPSAutocertServer(httpsAddr, s.newEngine(), m)
		s.httpsServer.ErrorLog = log.ToStdLogger(log.DefaultLogger, log.LevelWarn)

		s.httpServer = s.makeHTTPServer(httpAddr, s.httpChallengeRouter(m))
		s.httpServer.ErrorLog = log.ToStdLogger(log.DefaultLogger, log.LevelWarn)

		s.lock.Unlock()

		go func() {
			log.Infof("activate http challenge server on %s", httpAddr)

			err := s.httpServer.ListenAndServe()
			log.Warnf("http challenge server terminated, %s", err)
		}()

		err := s.httpsServer.ListenAndServeTLS("", "")
		log.Warnf("https server terminated, %s", err)
	}
}

// Shutdown rest http server
func (s *Server) Shutdown() {
	log.Warnf("shutdown rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Debugf("http shutdown error, %s", err)
		}
		log.Debugf("shutdown http server completed")
	}

	if s.httpsServer != nil {
		log.Warnf("shutdown https server")
		if err := s.httpsServer.Shutdown(ctx); err != nil {
			log.Debugf("https shutdown error, %s", err)
		}
		log.Debugf("shutdown https server completed")
	}
	s.lock.Unlock()
}

func (s *Server) makeHTTPServer(addr string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}

func (s *Server) newEngine() *gin.Engine {
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(middleware.Throttle(CuncurrentRequests), middleware.RequestID)
	if s.cfg.Server.AccessLog {
		engine.Use(middleware.Logger)
	}
	engine.Use(middleware.Recovery)
	engine.Use(ginCors.New(ginCors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Mode", "X-XSRF-Token", "X-JWT"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Authorization"},
		MaxAge:           300,
		AllowWebSockets:  false,
	}))
	engine.HTMLRender = &tplEngine{}

	authHandler, _ := s.authenticator.Handlers()
	authMiddleware := s.authenticator.Middleware()

	staticCtrl := newStaticController(s.cfg.Server.BasePath(), s.version, s.cfg.Auth.Anonymous.Enabled)
	fileCtrl := newFileController(s.fileBrowserSvc, s.hasher)

	engine.NoRoute(staticCtrl.indexHandler)
	router := engine.Group(s.cfg.Server.BasePath())
	router.GET("/static/*path", staticCtrl.staticHandler)
	router.Any("/auth/*path", middleware.NoCache, gin.WrapH(authHandler))

	v1 := router.Group("/api/v1")
	{
		protected := v1.Group("")
		{
			protected.Use(middleware.Timeout(ProtectedRoutesTimeout))
			protected.Use(middleware.LimitHandler(tollbooth.NewLimiter(StaticRouterLimiter, nil)))
			protected.Use(middleware.WrapHH(authMiddleware.Auth), middleware.User(), middleware.NoCache)

			// file handlers
			protected.GET("/files/:volume/*path", fileCtrl.ListHandler)
			//protected.DELETE("/files/:volume/*path", fileCtrl.DeleteHandler)
			//protected.GET("/files/volumes/:id/*path", fileCtrl.VolumeListHandler)
		}
	}

	/*v1 := router.Group("/api/v1")
	{
		public := v1.Group("")
		{
			public.Use(middleware.Timeout(PublicRoutesTimeout))
			public.Use(middleware.LimitHandler(tollbooth.NewLimiter(StaticRouterLimiter, nil)))
			public.GET("/avatar/*path", gin.WrapH(avatarHandler))
		}

		protected := v1.Group("")
		{
			protected.Use(middleware.Timeout(ProtectedRoutesTimeout))
			protected.Use(middleware.LimitHandler(tollbooth.NewLimiter(StaticRouterLimiter, nil)))
			protected.Use(middleware.WrapHH(authMiddleware.Auth), middleware.User(s.Store), middleware.NoCache)
			protected.GET("/resources/*path", fileCtrl.ListHandler)
			protected.POST("/resources/*path", fileCtrl.ModifyHandler)
			protected.PUT("/resources/*path", fileCtrl.ModifyHandler)
			protected.DELETE("/resources/*path", fileCtrl.DeleteHandler)
		}
	}*/

	return engine
}
