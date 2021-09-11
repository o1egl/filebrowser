//go:generate go-enum --marshal --nocase --names --file $GOFILE
package config

import (
	"net/url"
	"time"

	"github.com/google/uuid"
	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	RootPath string `yaml:"root"` // server root folder
	Secret   string `yaml:"secret" validate:"required"`
	Locale   string `yaml:"locale"` // default locale
	Server   Server `yaml:"server"` // http server config
	Auth     Auth   `yaml:"auth"`
	Store    Store  `yaml:"store"`
}

func (c *Config) Validate() error {
	return validator.New().Struct(c)
}

type Server struct {
	AccessLog bool   `yaml:"access_log"`
	URL       string `yaml:"url"` // file browser url. required for ssl and oauth
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	SSL       SSL    `yaml:"ssl"`
}

// BasePath returns base path for the server.
// For example for serverURL https://filebrowser.org/base/path it should return /base/path
func (s Server) BasePath() string {
	u, err := url.Parse(s.URL)
	if err != nil {
		return "/"
	}
	return u.Path
}

// Hostname returns hostname for the server.
// For example for serverURL https://filebrowser.org:443 it should return filebrowser.org
func (s Server) Hostname() string {
	u, err := url.Parse(s.URL)
	if err != nil {
		return ""
	}
	return u.Hostname()
}

type Auth struct {
	TTL struct {
		JWT    time.Duration `yaml:"jwt"`    // jwt TTL
		Cookie time.Duration `yaml:"cookie"` // auth cookie TTL
	} `yaml:"ttl"`
	Google    OAuth     `yaml:"google"`   // google oauth
	Github    OAuth     `yaml:"github"`   // github oauth
	Facebook  OAuth     `yaml:"facebook"` // facebook oauth
	Twitter   OAuth     `yaml:"twitter"`  // twitter oauth
	Dev       bool      `yaml:"dev"`      // enable dev (local) oauth2
	User      User      `yaml:"user"`
	Anonymous Anonymous `yaml:"anonymous"`
}

type OAuth struct {
	CID  string `yaml:"cid"`  // OAuth client ID
	CSEC string `yaml:"csec"` // OAuth client secret
}

type User struct {
	Home        string      `yaml:"home"`        // default user home
	Permissions Permissions `yaml:"permissions"` // default user permissions for home volume
}
type Anonymous struct {
	Enabled     bool        `yaml:"enabled"`
	Home        string      `yaml:"home"`        // home path
	Permissions Permissions `yaml:"permissions"` // default user permissions for home volume
}

type Permissions struct {
	Admin    bool `yaml:"admin"`
	Create   bool `yaml:"create"`
	Rename   bool `yaml:"rename"`
	Modify   bool `yaml:"modify"`
	Delete   bool `yaml:"delete"`
	Share    bool `yaml:"share"`
	Download bool `yaml:"download"`
}

/*
ENUM(
sqlite
mysql
postgres
)
*/
type StoreType int

type Store struct {
	Type     StoreType     `yaml:"type"` // storage backend type
	SQLite   SQLiteStore   `yaml:"sqlite"`
	Mysql    MysqlStore    `yaml:"mysql"`
	Postgres PostgresStore `yaml:"postgres"`
}

type SQLiteStore struct {
	File string `yaml:"file"` // sqlite file location
}

type MysqlStore struct {
	DSN string `yaml:"dsn"` // mysql dsn (username:password@protocol(address)/dbname)
}

type PostgresStore struct {
	DSN string `yaml:"dsn"` // postgres dsn (postgres://username:password@address/dbname?sslmode=disable)
}

/*
ENUM(
none
static
auto
)
*/
type SSLMode int

type SSL struct {
	Mode SSLMode `yaml:"mode"`
	Port int     `yaml:"port"` // port number for https server
	Cert string  `yaml:"cert"` // path to cert.pem file
	Key  string  `yaml:"key"`  // path to key.pem file
	ACME ACME    `yaml:"acme"` // acme config
}

type ACME struct {
	Path  string   `yaml:"path"`  // dir where certificates will be stored by autocert manager
	Email string   `yaml:"email"` // admin email for certificate notifications
	FQDNs []string `yaml:"fqdns"` // FQDN(s) for ACME certificates
}

// Default returns default config
func Default() *Config {
	return &Config{
		RootPath: ".",
		Secret:   uuid.New().String(),
		Locale:   "en",
		Server: Server{
			Host: "0.0.0.0",
			Port: 8080,
			SSL: SSL{
				Mode: SSLModeNone,
				Port: 8443,
				ACME: ACME{
					Path: "./var/acme",
				},
			},
		},
		Auth: Auth{
			TTL: struct {
				JWT    time.Duration `yaml:"jwt"`
				Cookie time.Duration `yaml:"cookie"`
			}{
				JWT:    10 * time.Second,
				Cookie: 7 * 24 * time.Hour,
			},
			User: User{
				Home: "/",
				Permissions: Permissions{
					Admin:    false,
					Create:   true,
					Rename:   true,
					Modify:   true,
					Delete:   true,
					Share:    true,
					Download: true,
				},
			},
			Anonymous: Anonymous{
				Enabled: false,
				Home:    "/",
				Permissions: Permissions{
					Admin:    false,
					Create:   false,
					Rename:   false,
					Modify:   false,
					Delete:   false,
					Share:    false,
					Download: true,
				},
			},
		},
		Store: Store{
			Type: StoreTypeSqlite,
			SQLite: SQLiteStore{
				File: "./var/filebrowser.db",
			},
		},
	}
}
