//go:generate ../../tools/bin/go-enum --marshal --nocase --names --file $GOFILE
package config

import (
	"net/url"
	"time"
)

type Config struct {
	Root   string `yaml:"root"` // server root folder
	Secret string `yaml:"secret"`
	Locale string `yaml:"locale"` // default locale
	Upload Upload `yaml:"upload"`
	Server Server `yaml:"server"` // http server config
	Auth   Auth   `yaml:"auth"`
	Store  Store  `yaml:"store"`
}

func (c *Config) Validate() error {
	return nil
}

type Upload struct {
	Path      string `yaml:"path"`
	ChunkSize int64  `yaml:"chunk_size"`
	MaxSize   int64  `yaml:"max_size"`
}

type Server struct {
	AccessLog bool   `yaml:"access_log"`
	URL       string `yaml:"url"` // file browser url. required for ssl and oauth
	Bind      string `yaml:"bind"`
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

type AuthTTL struct {
	JWT    time.Duration `yaml:"jwt"`    // jwt TTL
	Cookie time.Duration `yaml:"cookie"` // auth cookie TTL
}

type Auth struct {
	TTL AuthTTL `yaml:"ttl"`
}

type VolumePermissions struct {
	Read   bool `yaml:"read"`
	Create bool `yaml:"create"`
	Modify bool `yaml:"modify"`
	Delete bool `yaml:"delete"`
	Share  bool `yaml:"share"`
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
	Bind string  `yaml:"port"` // server addr
	Cert string  `yaml:"cert"` // path to cert.pem file
	Key  string  `yaml:"key"`  // path to key.pem file
	ACME ACME    `yaml:"acme"` // acme config
}

type ACME struct {
	Path  string   `yaml:"path"`  // dir where certificates will be stored by autocert manager
	Email string   `yaml:"email"` // admin email for certificate notifications
	FQDNs []string `yaml:"fqdns"` // FQDN(s) for ACME certificates
}
