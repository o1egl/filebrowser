package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/didip/tollbooth/v6"
	ginCors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/lcw"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v3/backend/log"
	"github.com/filebrowser/filebrowser/v3/backend/rest/middleware"
	"github.com/filebrowser/filebrowser/v3/backend/token"
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
	Root          afero.Fs
	Authenticator *auth.Service
	TokenService  *token.Service
	Store         Store
	Cache         LoadingCache
	Host          string
	Port          int
	ServerURL     string
	SharedSecret  string
	Revision      string
	AccessLog     bool
	Anonymous     bool

	SSLConfig   SSLConfig
	httpsServer *http.Server
	httpServer  *http.Server
	lock        sync.Mutex
}

type Store interface {
	middleware.UserStore
}

// LoadingCache defines interface for caching
type LoadingCache interface {
	Get(key lcw.Key, fn func() ([]byte, error)) (data []byte, err error) // load from cache if found or put to cache and return
	Flush(req lcw.FlusherRequest)                                        // evict matched records
}

func (s *Server) Run() {
	httpAddr := fmt.Sprintf("%s:%d", s.Host, s.Port)
	httpsAddr := fmt.Sprintf("%s:%d", s.Host, s.SSLConfig.Port)
	switch s.SSLConfig.SSLMode {
	case None:
		log.Infof("activate http rest server on %s", httpAddr)

		s.lock.Lock()
		s.httpServer = s.makeHTTPServer(httpAddr, s.newEngine())
		s.httpServer.ErrorLog = log.ToStdLogger(log.DefaultLogger, log.LevelWarn)
		s.lock.Unlock()

		err := s.httpServer.ListenAndServe()
		log.Warnf("http server terminated, %s", err)
	case Static:
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

		err := s.httpsServer.ListenAndServeTLS(s.SSLConfig.Cert, s.SSLConfig.Key)
		log.Warnf("https server terminated, %s", err)
	case Auto:
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
	if s.AccessLog {
		engine.Use(middleware.Logger)
	}
	engine.Use(middleware.Recovery)
	engine.Use(ginCors.New(ginCors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-XSRF-Token", "X-JWT"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Authorization"},
		MaxAge:           300,
		AllowWebSockets:  true,
	}))
	engine.HTMLRender = &tplEngine{Reload: true}

	authHandler, avatarHandler := s.Authenticator.Handlers()
	authMiddleware := s.Authenticator.Middleware()

	publicCtrl, fileCtrl := s.makeHandlerGroups()

	engine.NoRoute(publicCtrl.indexHandler)
	router := engine.Group(s.getServerBasePath())
	router.GET("/static/*path", publicCtrl.staticHandler)
	router.Any("/auth/*path", middleware.NoCache, gin.WrapH(authHandler))
	v1 := router.Group("/api/v1")
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
	}

	return engine
}

func (s *Server) makeHandlerGroups() (*publicHandlers, *fileController) {
	publicHandlers := &publicHandlers{
		BasePath:  s.getServerBasePath(),
		Revision:  s.Revision,
		Anonymous: s.Anonymous,
	}

	fileCtrl := &fileController{
		root: s.Root,
	}

	return publicHandlers, fileCtrl
}

// getServerBasePath returns base path for the server.
// For example for serverURL https://filebrowser.org/base/path it should return /base/path
func (s *Server) getServerBasePath() string {
	u, err := url.Parse(s.ServerURL)
	if err != nil {
		return "/"
	}
	return u.Path
}
