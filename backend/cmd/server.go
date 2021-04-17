package cmd

import (
	"context"
	"crypto/sha1" //nolint:gosec
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/provider"
	authToken "github.com/go-pkgz/auth/token"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v3/store/sql"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"

	"github.com/filebrowser/filebrowser/v3/cache"
	"github.com/filebrowser/filebrowser/v3/hash"
	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/rest/api"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/token"
)

// ServerCommand with command line flags and env
type ServerCommand struct {
	Auth      AuthGroup      `group:"auth" namespace:"auth" env-namespace:"AUTH"`
	Cache     CacheGroup     `group:"cache" namespace:"cache" env-namespace:"CACHE"`
	Store     StoreGroup     `group:"store" namespace:"store" env-namespace:"STORE"`
	SSL       SSLGroup       `group:"ssl" namespace:"ssl" env-namespace:"SSL"`
	AccessLog AccessLogGroup `group:"access-log" namespace:"access-log" env-namespace:"ACCESS_LOG"`

	ServerURL string `long:"url" env:"SERVER_URL" description:"file browser url"`
	Secret    string `long:"secret" env:"SECRET" description:"shared secret key"`
	Locale    string `long:"locale" env:"LOCALE" default:"en" description:"default locale"`
	RootPath  string `long:"root" env:"ROOT_PATH" default:"." description:"root folder"`
	Host      string `long:"host" env:"HOST" default:"0.0.0.0" description:"host"`
	Port      int    `long:"port" env:"PORT" default:"8080" description:"port"`

	CommonOpts
}

// StoreGroup defines options group for storage
type StoreGroup struct {
	Type   string `long:"type" env:"TYPE" description:"type of storage" choice:"sqlite" choice:"postgres" choice:"mysql" default:"sqlite"`
	SQLite struct {
		File string `long:"file" env:"FILE" default:"./var/filebrowser.db" description:"sqlite file location"`
	} `group:"sqlite" namespace:"sqlite" env-namespace:"SQLITE"`
	Postgres struct {
		DSN string `long:"dsn" env:"DSN" description:"postgres dsn (postgres://username:password@address/dbname?sslmode=disable)"`
	} `group:"postgres" namespace:"postgres" env-namespace:"POSTGRES"`
	MySQL struct {
		DSN string `long:"dsn" env:"DSN" description:"mysql dsn (username:password@protocol(address)/dbname)"`
	} `group:"mysql" namespace:"mysql" env-namespace:"MYSQL"`
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
			Create   bool `long:"create" env:"CREATE" description:"add create permission"`
			Rename   bool `long:"rename" env:"RENAME" description:"add rename permission"`
			Modify   bool `long:"modify" env:"MODIFY" description:"add modify permission"`
			Delete   bool `long:"delete" env:"DELETE" description:"add delete permission"`
			Share    bool `long:"share" env:"SHARE" description:"add share permission"`
			Download bool `long:"download" env:"DOWNLOAD" description:"add download permission"`
		} `group:"perm" namespace:"perm" env-namespace:"PERM"`
	} `group:"anon" namespace:"anon" env-namespace:"ANON"`
	Admin struct {
		Password string `long:"password" env:"PASSWORD" default:"admin" description:"encrypted admin password"`
	} `group:"admin" namespace:"admin" env-namespace:"ADMIN"`
}

// OAuthGroup defines options group for oauth params
type OAuthGroup struct {
	CID  string `long:"cid" env:"CID" description:"OAuth client ID"`
	CSEC string `long:"csec" env:"CSEC" description:"OAuth client secret"`
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
	Type         string   `long:"type" env:"TYPE" description:"ssl (auto) support" choice:"none" choice:"static" choice:"auto" default:"none"` //nolint
	Port         int      `long:"port" env:"PORT" description:"port number for https server" default:"8443"`
	Cert         string   `long:"cert" env:"CERT" description:"path to cert.pem file"`
	Key          string   `long:"key" env:"KEY" description:"path to key.pem file"`
	ACMELocation string   `long:"acme-location" env:"ACME_LOCATION" description:"dir where certificates will be stored by autocert manager" default:"./var/acme"` //nolint
	ACMEEmail    string   `long:"acme-email" env:"ACME_EMAIL" description:"admin email for certificate notifications"`
	FQDNs        []string `long:"fqdn" env:"ACME_FQDN" env-delim:"," description:"FQDN(s) for ACME certificates"`
}

// AccessLogGroup defines options group for access log
type AccessLogGroup struct {
	Enable bool `long:"enable" env:"ENABLE" description:"enable access log"`
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
		"STORE_POSTGRES_DSN", "STORE_MYSQL_DSN",
	)

	ctx, cancel := context.WithCancel(context.Background())
	go func() { // catch signal and invoke graceful termination
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
		<-stop
		log.Warnf("interrupt signal")
		cancel()
	}()

	app, err := s.newServerApp(ctx)
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
func (s *ServerCommand) newServerApp(ctx context.Context) (*serverApp, error) {
	entClient, err := s.makeEntClient(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initiate data engine")
	}

	userStore := sql.NewUserStore(entClient)

	/*if err := s.createAdminUser(dataStore); err != nil {
		return nil, errors.Wrap(err, "failed to create admin user")
	}

	if err := s.createAnonymousUser(dataStore); err != nil {
		return nil, errors.Wrap(err, "failed to create anonymous user")
	}*/

	authRefreshCache := newAuthRefreshCache()
	localAuthProvider := newLocalAuthProvider(userStore)
	authenticator, err := s.makeAuthenticator(userStore, authRefreshCache, localAuthProvider)
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
		TokenService:  token.New(s.Secret),
		UserStore:     userStore,
		Host:          s.Host,
		Port:          s.Port,
		ServerURL:     s.ServerURL,
		SharedSecret:  s.Secret,
		Revision:      s.Revision,
		AccessLog:     s.AccessLog.Enable,
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
	dataStore store.UserStore,
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
			return s.Secret, nil
		}),
		ClaimsUpd: s.newClaimsUpdater(context.Background(), dataStore),
		BasicAuthChecker: func(user, passwd string) (ok bool, userInfo authToken.User, err error) {
			u, err := dataStore.GetByUsernameAndProvider(context.Background(), user, "local")
			if err != nil {
				return false, authToken.User{}, err
			}
			if !hash.CheckPassword(passwd, u.Password) {
				return false, authToken.User{}, errors.New("basic auth credentials check failed")
			}
			return true, authToken.User{
				Name: u.Name,
				ID:   u.ID,
			}, nil
		},
		Validator: authToken.ValidatorFunc(func(token string, claims authToken.Claims) bool { // check on each auth call (in middleware)
			if claims.User == nil {
				return false
			}
			return !claims.User.BoolAttr("blocked")
		}),
		AvatarStore:  avatar.NewNoOp(),
		Logger:       log.NewLogrAdapter(log.DefaultLogger),
		RefreshCache: authRefreshCache,
	})

	s.addAuthProviders(authenticator, localProvider)

	return authenticator, nil
}

func (s *ServerCommand) newClaimsUpdater(ctx context.Context, userStore store.UserStore) authToken.ClaimsUpdater {
	return authToken.ClaimsUpdFunc(func(c authToken.Claims) authToken.Claims {
		if c.User == nil {
			return c
		}
		user, err := userStore.Get(ctx, c.User.ID)
		switch {
		// new login with external provider
		case errors.Is(err, store.ErrNotFound):
			user = &store.User{
				ID:           c.User.ID,
				Name:         c.User.Name,
				Provider:     strings.Split(c.User.ID, "_")[0],
				Username:     "",
				Password:     "",
				Scope:        s.Auth.User.Scope,
				Locale:       s.Locale,
				LockPassword: false,
				Blocked:      false,
			}
			if err := userStore.Save(ctx, user); err != nil {
				log.WithContext(ctx).Errorf("failed to create user: %+v", err)
				return c
			}
		}
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

func newLocalAuthProvider(userStore store.UserStore) provider.CredChecker {
	return provider.CredCheckerFunc(func(username, password string) (ok bool, err error) {
		user, err := userStore.Get(context.TODO(), username)
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

func (s *ServerCommand) makeEntClient(ctx context.Context) (client *ent.Client, err error) {
	switch s.Store.Type {
	case "sqlite":
		if err = makeDirs(filepath.Dir(s.Store.SQLite.File)); err != nil {
			return nil, errors.Wrap(err, "failed to create sqlite store")
		}
		client, err = ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_fk=1", s.Store.SQLite.File))
	case "postgres":
		client, err = ent.Open("postgres", s.Store.Postgres.DSN)
	case "mysql":
		client, err = ent.Open("mysql", s.Store.Postgres.DSN)
	default:
		return nil, errors.Errorf("unsupported store type %s", s.Store.Type)
	}
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize data store")
	}

	log.WithContext(ctx).Debugf("Apply schema migrations")
	if err := client.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return client, nil
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

func (s *ServerCommand) createAdminUser(storage store.UserStore) error {
	pwdHash, err := hash.Password(s.Auth.Admin.Password)
	if err != nil {
		return err
	}
	user := &store.User{
		ID:           "local_" + authToken.HashID(sha1.New(), "admin"), //nolint:gosec
		Name:         "Admin",
		Provider:     "local",
		Username:     "admin",
		Password:     pwdHash,
		Scope:        "/",
		Locale:       s.Locale,
		LockPassword: false,
		Blocked:      false,
	}
	return storage.Save(context.Background(), user)
}

func (s *ServerCommand) createAnonymousUser(storage store.UserStore) error {
	pwdHash, err := hash.Password("")
	if err != nil {
		return err
	}
	user := &store.User{
		ID:           "local_" + authToken.HashID(sha1.New(), "anonymous"), //nolint:gosec
		Name:         "Anonymous",
		Provider:     "local",
		Username:     "anonymous",
		Password:     pwdHash,
		Scope:        s.Auth.Anonymous.Scope,
		Locale:       s.Locale,
		LockPassword: true,
		Blocked:      !s.Auth.Anonymous.Enable,
	}
	return storage.Save(context.Background(), user)
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
	cache.Cache
}

func newAuthRefreshCache() *authRefreshCache {
	memCache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000,
		MaxCost:     1000,
		BufferItems: 64,
	})

	if err != nil {
		log.Fatalf("Failed to init cache: %s", err)
	}

	return &authRefreshCache{Cache: memCache}
}

// Get implements cache getter with key converted to string
func (c *authRefreshCache) Get(key interface{}) (interface{}, bool) {
	return c.Get(key)
}

// Set implements cache setter with key converted to string
func (c *authRefreshCache) Set(key, value interface{}) {
	c.SetWithTTL(key, value, 1, 5*time.Second)
}
