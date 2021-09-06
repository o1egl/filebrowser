package server

import (
	"context" //nolint:gosec
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/filebrowser/filebrowser/v3/cache"
	"github.com/filebrowser/filebrowser/v3/cmd"
	"github.com/filebrowser/filebrowser/v3/log"
	pkgAuth "github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/provider"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
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

	cmd.CommonOpts
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
		Home        string `long:"home" env:"HOME" default:"/" description:"default user home. must start with /"`
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
		Home        string `long:"home" env:"HOME" default:"/" description:"user home. must start with /"`
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

// Execute runs file browser server
func (s *ServerCommand) Execute(_ []string) error {
	cmd.ResetEnv(
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

	app, err := InitializeServer(ctx, s)
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

func addAuthProviders(srvCmd *ServerCommand, authenticator *pkgAuth.Service, localProvider provider.CredChecker) {
	providers := 0

	providers++
	authenticator.AddDirectProvider("local", localProvider)

	if srvCmd.Auth.Google.CID != "" && srvCmd.Auth.Google.CSEC != "" {
		authenticator.AddProvider("google", srvCmd.Auth.Google.CID, srvCmd.Auth.Google.CSEC)
		providers++
	}
	if srvCmd.Auth.Github.CID != "" && srvCmd.Auth.Github.CSEC != "" {
		authenticator.AddProvider("github", srvCmd.Auth.Github.CID, srvCmd.Auth.Github.CSEC)
		providers++
	}
	if srvCmd.Auth.Facebook.CID != "" && srvCmd.Auth.Facebook.CSEC != "" {
		authenticator.AddProvider("facebook", srvCmd.Auth.Facebook.CID, srvCmd.Auth.Facebook.CSEC)
		providers++
	}
	if srvCmd.Auth.Twitter.CID != "" && srvCmd.Auth.Twitter.CSEC != "" {
		authenticator.AddProvider("twitter", srvCmd.Auth.Twitter.CID, srvCmd.Auth.Twitter.CSEC)
		providers++
	}

	if srvCmd.Auth.Dev {
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
