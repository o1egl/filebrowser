package cmd

import (
	"context"
	"crypto/sha1" //nolint:gosec
	"net/url"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	authToken "github.com/go-pkgz/auth/token"
	cache "github.com/go-pkgz/lcw"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	bolt "go.etcd.io/bbolt"

	"github.com/filebrowser/filebrowser/v3/hash"
	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/rest/api"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/engine"
	"github.com/filebrowser/filebrowser/v3/store/service"
	"github.com/filebrowser/filebrowser/v3/token"
)

// ServerCommand with command line flags and env
type ServerCommand struct {
	Auth   AuthGroup   `group:"auth" namespace:"auth" env-namespace:"AUTH"`
	Avatar AvatarGroup `group:"avatar" namespace:"avatar" env-namespace:"AVATAR"`
	Cache  CacheGroup  `group:"cache" namespace:"cache" env-namespace:"CACHE"`
	Store  StoreGroup  `group:"store" namespace:"store" env-namespace:"STORE"`
	SSL    SSLGroup    `group:"ssl" namespace:"ssl" env-namespace:"SSL"`

	Locale    string `long:"locale" env:"LOCALE" default:"en" description:"default locale"`
	AccessLog bool   `long:"enable-access-log" env:"ENABLE_ACCESS_LOG" description:"enable access log"`
	RootPath  string `long:"root" env:"ROOT_PATH" default:"." description:"root folder"`
	Host      string `long:"host" env:"HOST" default:"0.0.0.0" description:"host"`
	Port      int    `long:"port" env:"PORT" default:"8080" description:"port"`

	CommonOpts
}

// AuthGroup defines options group store params
type StoreGroup struct {
	Type string `long:"type" env:"TYPE" description:"type of storage" choice:"bolt" default:"bolt"`
	Bolt struct {
		File    string        `long:"file" env:"FILE" default:"./var/filebrowser.db" description:"bolt file location"`
		Timeout time.Duration `long:"timeout" env:"TIMEOUT" default:"30s" description:"bolt timeout"`
	} `group:"bolt" namespace:"bolt" env-namespace:"BOLT"`
}

// AuthGroup defines options group auth params
type AuthGroup struct {
	TTL struct {
		JWT    time.Duration `long:"jwt" env:"JWT" default:"10s" description:"jwt TTL"`
		Cookie time.Duration `long:"cookie" env:"COOKIE" default:"200h" description:"auth cookie TTL"`
	} `group:"ttl" namespace:"ttl" env-namespace:"TTL"`
	Google   OAuthGroup `group:"google" namespace:"google" env-namespace:"GOOGLE" description:"Google OAuth"`
	Github   OAuthGroup `group:"github" namespace:"github" env-namespace:"GITHUB" description:"Github OAuth"`
	Facebook OAuthGroup `group:"facebook" namespace:"facebook" env-namespace:"FACEBOOK" description:"Facebook OAuth"`
	Twitter  OAuthGroup `group:"twitter" namespace:"twitter" env-namespace:"TWITTER" description:"Twitter OAuth"`
	Dev      bool       `long:"dev" env:"DEV" description:"enable dev (local) oauth2"`
	User     struct {
		Scope       string `long:"scope" env:"SCOPE" default:"/" description:"default user scope. must start with /"`
		Permissions struct {
			Execute  bool `long:"execute" env:"EXECUTE" description:"add execute permission"`
			Create   bool `long:"create" env:"CREATE" description:"add create permission"`
			Rename   bool `long:"rename" env:"RENAME" description:"add rename permission"`
			Modify   bool `long:"modify" env:"MODIFY" description:"add modify permission"`
			Delete   bool `long:"delete" env:"DELETE" description:"add delete permission"`
			Share    bool `long:"share" env:"SHARE" description:"add share permission"`
			Download bool `long:"download" env:"DOWNLOAD" description:"add download permission"`
		} `group:"perm" namespace:"perm" env-namespace:"PERM"`
	} `group:"user" namespace:"user" env-namespace:"USER"`
	Anonymous struct {
		Enable      bool   `long:"enable" env:"ENABLE" description:"enable anonymous user"`
		Scope       string `long:"scope" env:"SCOPE" default:"/" description:"user scope. must start with /"`
		Permissions struct {
			Admin    bool `long:"admin" env:"ADMIN" description:"add admin permission"`
			Execute  bool `long:"execute" env:"EXECUTE" description:"add execute permission"`
			Create   bool `long:"create" env:"CREATE" description:"add create permission"`
			Rename   bool `long:"rename" env:"RENAME" description:"add rename permission"`
			Modify   bool `long:"modify" env:"MODIFY" description:"add modify permission"`
			Delete   bool `long:"delete" env:"DELETE" description:"add delete permission"`
			Share    bool `long:"share" env:"SHARE" description:"add share permission"`
			Download bool `long:"download" env:"DOWNLOAD" description:"add download permission"`
		} `group:"perm" namespace:"perm" env-namespace:"PERM"`
	} `group:"anon" namespace:"anon" env-namespace:"ANON"`
	Admin struct {
		Username string `long:"username" env:"USERNAME" default:"admin" description:"admin username"`
		Password string `long:"password" env:"PASSWORD" default:"admin" description:"admin password"`
	} `group:"admin" namespace:"admin" env-namespace:"ADMIN"`
}

// OAuthGroup defines options group for oauth params
type OAuthGroup struct {
	CID  string `long:"cid" env:"CID" description:"OAuth client ID"`
	CSEC string `long:"csec" env:"CSEC" description:"OAuth client secret"`
}

// AvatarGroup defines options group for avatar params
type AvatarGroup struct {
	Type string `long:"type" env:"TYPE" description:"type of avatar storage" choice:"fs" choice:"bolt" choice:"uri" default:"fs"` //nolint
	FS   struct {
		Path string `long:"path" env:"PATH" default:"./var/avatars" description:"avatars location"`
	} `group:"fs" namespace:"fs" env-namespace:"FS"`
	Bolt struct {
		File string `long:"file" env:"FILE" default:"./var/avatars.db" description:"avatars bolt file location"`
	} `group:"bolt" namespace:"bolt" env-namespace:"bolt"`
	URI    string `long:"uri" env:"URI" default:"./var/avatars" description:"avatar's store URI"`
	RszLmt int    `long:"rsz-lmt" env:"RESIZE" default:"0" description:"max image size for resizing avatars on save"`
}

// CacheGroup defines options group for cache params
type CacheGroup struct {
	Type string `long:"type" env:"TYPE" description:"type of cache" choice:"mem" choice:"none" default:"mem"` // nolint
	Max  struct {
		Items int   `long:"items" env:"ITEMS" default:"1000" description:"max cached items"`
		Value int   `long:"value" env:"VALUE" default:"65536" description:"max size of cached value"`
		Size  int64 `long:"size" env:"SIZE" default:"50000000" description:"max size of total cache"`
	} `group:"max" namespace:"max" env-namespace:"MAX"`
}

// SSLGroup defines options group for server ssl params
type SSLGroup struct {
	Type         string `long:"type" env:"TYPE" description:"ssl (auto) support" choice:"none" choice:"static" choice:"auto" default:"none"` //nolint
	Port         int    `long:"port" env:"PORT" description:"port number for https server" default:"8443"`
	Cert         string `long:"cert" env:"CERT" description:"path to cert.pem file"`
	Key          string `long:"key" env:"KEY" description:"path to key.pem file"`
	ACMELocation string `long:"acme-location" env:"ACME_LOCATION" description:"dir where certificates will be stored by autocert manager" default:"./var/acme"` //nolint
	ACMEEmail    string `long:"acme-email" env:"ACME_EMAIL" description:"admin email for certificate notifications"`
}

// LoadingCache defines interface for caching
type LoadingCache interface {
	Get(key cache.Key, fn func() ([]byte, error)) (data []byte, err error) // load from cache if found or put to cache and return
	Flush(req cache.FlusherRequest)                                        // evict matched records
	Close() error
}

// serverApp holds all active objects
type serverApp struct {
	*ServerCommand
	restSrv    *api.Server
	terminated chan struct{}
}

// Execute runs file browser server
func (s *ServerCommand) Execute(_ []string) error {
	resetEnv(
		"SECRET", "AUTH_ADMIN_USERNAME", "AUTH_ADMIN_PASSWORD",
		"GOOGLE_CID", "GOOGLE_CSEC",
		"GITHUB_CID", "GITHUB_CSEC",
		"FACEBOOK_CID", "FACEBOOK_CSEC",
		"TWITTER_CID", "TWITTER_CSEC",
	)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warnf("interrupt signal")
		cancel()
	}()

	app, err := s.newServerApp()
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

// newServerApp prepares application and return it with all active parts
// doesn't start anything
func (s *ServerCommand) newServerApp() (*serverApp, error) {
	loadingCache, err := s.makeCache()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make cache")
	}

	dataEngine, err := s.makeDataEngine()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initiate data engine")
	}

	dataStore := &service.DataStore{
		Engine: dataEngine,
	}

	if err := s.createAdminUser(dataStore); err != nil {
		return nil, errors.Wrap(err, "failed to create admin user")
	}

	if err := s.createAnonymousUser(dataStore); err != nil {
		return nil, errors.Wrap(err, "failed to create anonymous user")
	}

	avatarStore, err := s.makeAvatarStore()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make avatar store")
	}

	authRefreshCache := newAuthRefreshCache()
	localAuthProvider := newLocalAuthProvider(dataStore)
	authenticator, err := s.makeAuthenticator(dataStore, avatarStore, authRefreshCache, localAuthProvider)
	if err != nil {
		return nil, errors.Wrap(err, "failed to make authenticator")
	}

	sslConfig, err := s.makeSSLConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to make config of ssl server params")
	}

	absRootPath, err := filepath.Abs(s.RootPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get abs path")
	}

	apiServer := &api.Server{
		Root:          afero.NewBasePathFs(afero.NewOsFs(), absRootPath),
		Authenticator: authenticator,
		TokenService:  token.New(s.SharedSecret),
		Store:         dataStore,
		Cache:         loadingCache,
		Host:          s.Host,
		Port:          s.Port,
		ServerURL:     s.ServerURL,
		SharedSecret:  s.SharedSecret,
		Revision:      s.Revision,
		AccessLog:     s.AccessLog,
		Anonymous:     s.Auth.Anonymous.Enable,
		SSLConfig:     sslConfig,
	}

	return &serverApp{
		ServerCommand: s,
		restSrv:       apiServer,
		terminated:    make(chan struct{}),
	}, nil
}

func (s *ServerCommand) makeAuthenticator(
	dataStore *service.DataStore,
	avatarStore avatar.Store,
	authRefreshCache *authRefreshCache, //nolint:interfacer
	localProvider provider.CredChecker,
) (*auth.Service, error) { //nolint:unparam
	authenticator := auth.NewService(auth.Opts{
		DisableXSRF:    true, // TODO remove it
		URL:            strings.TrimSuffix(s.ServerURL, "/"),
		Issuer:         "File Browser",
		TokenDuration:  s.Auth.TTL.JWT,
		CookieDuration: s.Auth.TTL.Cookie,
		SecureCookies:  strings.HasPrefix(s.ServerURL, "https://"),
		SecretReader: authToken.SecretFunc(func(aud string) (string, error) {
			if s.SharedSecret == "" {
				return "", errors.New("shared secret is not provided")
			}
			return s.SharedSecret, nil
		}),
		ClaimsUpd: s.newClaimsUpdater(context.Background(), dataStore),
		Validator: authToken.ValidatorFunc(func(token string, claims authToken.Claims) bool { // check on each auth call (in middleware)
			if claims.User == nil {
				return false
			}
			return !claims.User.BoolAttr("blocked")
		}),
		AvatarResizeLimit: s.Avatar.RszLmt,
		AvatarRoutePath:   path.Join(s.getServerBasePath(), "/api/v1/avatar"),
		AvatarStore:       avatarStore,
		Logger:            log.NewLogrAdapter(log.DefaultLogger),
		RefreshCache:      authRefreshCache,
		UseGravatar:       true,
	})

	s.addAuthProviders(authenticator, localProvider)

	return authenticator, nil
}

func (s *ServerCommand) newClaimsUpdater(ctx context.Context, dataStore *service.DataStore) authToken.ClaimsUpdater {
	return authToken.ClaimsUpdFunc(func(c authToken.Claims) authToken.Claims {
		if c.User == nil {
			return c
		}
		user, err := dataStore.FindUserByID(ctx, c.User.ID)
		switch {
		case errors.Is(err, store.ErrNotFound):
			user = &store.User{
				ID:           c.User.ID,
				Name:         c.User.Name,
				Picture:      c.User.Picture,
				Provider:     strings.Split(c.User.ID, "_")[0],
				Username:     "",
				Password:     "",
				Scope:        s.Auth.User.Scope,
				Locale:       s.Locale,
				Rules:        nil,
				Commands:     nil,
				LockPassword: false,
				Permissions: store.Permissions{
					Admin:    false,
					Execute:  s.Auth.User.Permissions.Execute,
					Create:   s.Auth.User.Permissions.Create,
					Rename:   s.Auth.User.Permissions.Rename,
					Modify:   s.Auth.User.Permissions.Modify,
					Delete:   s.Auth.User.Permissions.Delete,
					Share:    s.Auth.User.Permissions.Share,
					Download: s.Auth.User.Permissions.Download,
				},
				Blocked: false,
			}
			if err := dataStore.SaveUser(ctx, user); err != nil {
				log.WithContext(ctx).Errorf("failed to create user: %+v", err)
				return c
			}
		case err != nil:
			log.WithContext(ctx).Errorf("failed to find user: %+v", err)
			return c
		}
		c.User.Name = user.Name
		c.User.SetAdmin(user.Permissions.Admin)
		c.User.SetBoolAttr("anonymous", user.Username == "anonymous")
		c.User.SetBoolAttr("lockPassword", user.LockPassword)
		c.User.SetBoolAttr("permAdmin", user.Permissions.Admin)
		c.User.SetBoolAttr("permExecute", user.Permissions.CanExecute())
		c.User.SetBoolAttr("permCreate", user.Permissions.CanCreate())
		c.User.SetBoolAttr("permRename", user.Permissions.CanRename())
		c.User.SetBoolAttr("permModify", user.Permissions.CanModify())
		c.User.SetBoolAttr("permDelete", user.Permissions.CanDelete())
		c.User.SetBoolAttr("permShare", user.Permissions.CanShare())
		c.User.SetBoolAttr("permDownload", user.Permissions.CanDownload())
		c.User.SetStrAttr("viewMode", "mosaic")
		c.User.SetStrAttr("locale", user.Locale)
		c.User.SetBoolAttr("blocked", user.Blocked)

		return c
	})
}

func (s *ServerCommand) addAuthProviders(authenticator *auth.Service, localProvider provider.CredChecker) {
	providers := 0

	providers++
	authenticator.AddDirectProvider("local", localProvider)

	if s.Auth.Google.CID != "" && s.Auth.Google.CSEC != "" {
		authenticator.AddProvider("google", s.Auth.Google.CID, s.Auth.Google.CSEC)
		providers++
	}
	if s.Auth.Github.CID != "" && s.Auth.Github.CSEC != "" {
		authenticator.AddProvider("github", s.Auth.Github.CID, s.Auth.Github.CSEC)
		providers++
	}
	if s.Auth.Facebook.CID != "" && s.Auth.Facebook.CSEC != "" {
		authenticator.AddProvider("facebook", s.Auth.Facebook.CID, s.Auth.Facebook.CSEC)
		providers++
	}
	if s.Auth.Twitter.CID != "" && s.Auth.Twitter.CSEC != "" {
		authenticator.AddProvider("twitter", s.Auth.Twitter.CID, s.Auth.Twitter.CSEC)
		providers++
	}

	if s.Auth.Dev {
		log.Warnf("dev oauth provider is enabled")
		authenticator.AddProvider("dev", "", "")
		providers++

		// run dev/test oauth2 server on :8084
		go func() {
			devAuthServer, err := authenticator.DevAuth() // peak dev oauth2 server
			if err != nil {
				log.Fatalf("failed to start dev oauth2 server, %v", err)
			}
			devAuthServer.Run(context.Background())
		}()
	}

	if providers == 0 {
		log.Warnf("no auth providers defined")
	}
}

func newLocalAuthProvider(dataStore *service.DataStore) provider.CredChecker {
	return provider.CredCheckerFunc(func(username, password string) (ok bool, err error) {
		user, err := dataStore.FindUserByUsername(context.TODO(), username)
		if errors.Is(err, store.ErrNotFound) {
			return false, errors.New("user not found")
		}
		if err != nil {
			return false, err
		}

		if !hash.CheckPassword(password, user.Password) {
			return false, errors.New("incorrect user credentials")
		}
		return true, nil
	})
}

func (s *ServerCommand) makeDataEngine() (result engine.Interface, err error) {
	log.Infof("make data store, type=%s", s.Store.Type)
	switch s.Store.Type {
	case "bolt": //nolint:goconst
		if err = makeDirs(path.Dir(s.Store.Bolt.File)); err != nil {
			return nil, errors.Wrap(err, "failed to create bolt store")
		}
		result, err = engine.NewBoltDB(context.Background(), s.Store.Bolt.File, &bolt.Options{Timeout: s.Store.Bolt.Timeout})
	default:
		return nil, errors.Errorf("unsupported store type %s", s.Store.Type)
	}
	return result, errors.Wrap(err, "can't initialize data store")
}

func (s *ServerCommand) makeAvatarStore() (avatar.Store, error) {
	log.Infof("make avatar store, type=%s", s.Avatar.Type)

	switch s.Avatar.Type {
	case "fs":
		if err := makeDirs(s.Avatar.FS.Path); err != nil {
			return nil, errors.Wrap(err, "failed to create avatar store")
		}
		return avatar.NewLocalFS(s.Avatar.FS.Path), nil
	case "bolt":
		if err := makeDirs(path.Dir(s.Avatar.Bolt.File)); err != nil {
			return nil, errors.Wrap(err, "failed to create avatar store")
		}
		return avatar.NewBoltDB(s.Avatar.Bolt.File, bolt.Options{})
	case "uri":
		return avatar.NewStore(s.Avatar.URI)
	}
	return nil, errors.Errorf("unsupported avatar store type %s", s.Avatar.Type)
}

func (s *ServerCommand) makeCache() (LoadingCache, error) {
	log.Infof("make cache, type=%s", s.Cache.Type)
	switch s.Cache.Type {
	case "mem":
		backend, err := cache.NewLruCache(cache.MaxCacheSize(s.Cache.Max.Size), cache.MaxValSize(s.Cache.Max.Value),
			cache.MaxKeys(s.Cache.Max.Items))
		if err != nil {
			return nil, errors.Wrap(err, "cache backend initialization")
		}
		return cache.NewScache(backend), nil
	case "none": //nolint:goconst
		return cache.NewScache(&cache.Nop{}), nil
	}
	return nil, errors.Errorf("unsupported cache type %s", s.Cache.Type)
}

func (s *ServerCommand) makeSSLConfig() (config api.SSLConfig, err error) {
	switch s.SSL.Type {
	case "none":
		config.SSLMode = api.None
	case "static":
		if s.SSL.Cert == "" {
			return config, errors.New("path to cert.pem is required")
		}
		if s.SSL.Key == "" {
			return config, errors.New("path to key.pem is required")
		}
		config.SSLMode = api.Static
		config.Port = s.SSL.Port
		config.Cert = s.SSL.Cert
		config.Key = s.SSL.Key
	case "auto":
		config.SSLMode = api.Auto
		config.Port = s.SSL.Port
		config.ACMELocation = s.SSL.ACMELocation
		if s.SSL.ACMEEmail != "" {
			config.ACMEEmail = s.SSL.ACMEEmail
		} else if u, e := url.Parse(s.ServerURL); e == nil {
			config.ACMEEmail = "admin@" + u.Hostname()
		}
	}
	return config, err
}

func (s *ServerCommand) createAdminUser(storage *service.DataStore) error {
	pwdHash, err := hash.Password(s.Auth.Admin.Password)
	if err != nil {
		return err
	}
	user := &store.User{
		ID:           "local_" + authToken.HashID(sha1.New(), s.Auth.Admin.Username), //nolint:gosec
		Name:         "Admin",
		Picture:      "",
		Provider:     "local",
		Username:     s.Auth.Admin.Username,
		Password:     pwdHash,
		Scope:        "/",
		Locale:       s.Locale,
		Rules:        nil,
		Commands:     nil,
		LockPassword: true,
		Permissions: store.Permissions{
			Admin: true,
		},
		Blocked: false,
	}
	return storage.SaveUser(context.Background(), user)
}

func (s *ServerCommand) createAnonymousUser(storage *service.DataStore) error {
	pwdHash, err := hash.Password("")
	if err != nil {
		return err
	}
	user := &store.User{
		ID:           "local_" + authToken.HashID(sha1.New(), "anonymous"), //nolint:gosec
		Name:         "Anonymous",
		Picture:      "",
		Provider:     "local",
		Username:     "anonymous",
		Password:     pwdHash,
		Scope:        s.Auth.Anonymous.Scope,
		Locale:       s.Locale,
		Rules:        nil,
		Commands:     nil,
		LockPassword: true,
		Permissions: store.Permissions{
			Admin:    s.Auth.Anonymous.Permissions.Admin,
			Execute:  s.Auth.Anonymous.Permissions.Execute,
			Create:   s.Auth.Anonymous.Permissions.Create,
			Rename:   s.Auth.Anonymous.Permissions.Rename,
			Modify:   s.Auth.Anonymous.Permissions.Modify,
			Delete:   s.Auth.Anonymous.Permissions.Delete,
			Share:    s.Auth.Anonymous.Permissions.Share,
			Download: s.Auth.Anonymous.Permissions.Download,
		},
		Blocked: !s.Auth.Anonymous.Enable,
	}
	return storage.SaveUser(context.Background(), user)
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

// getServerBasePath returns base path for the server.
// For example for serverURL https://filebrowser.org/base/path it should return /base/path
func (s *ServerCommand) getServerBasePath() string {
	u, err := url.Parse(s.ServerURL)
	if err != nil {
		return "/"
	}
	return u.Path
}

// authRefreshCache used by authenticator to minimize repeatable token refreshes
type authRefreshCache struct {
	cache.LoadingCache
}

func newAuthRefreshCache() *authRefreshCache {
	expirableCache, _ := cache.NewExpirableCache(cache.TTL(5 * time.Minute)) //nolint:gomnd
	return &authRefreshCache{LoadingCache: expirableCache}
}

// Get implements cache getter with key converted to string
func (c *authRefreshCache) Get(key interface{}) (interface{}, bool) {
	return c.LoadingCache.Peek(key.(string))
}

// Set implements cache setter with key converted to string
func (c *authRefreshCache) Set(key, value interface{}) {
	_, _ = c.LoadingCache.Get(key.(string), func() (cache.Value, error) { return value, nil })
}
