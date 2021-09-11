//go:generate go-enum --marshal --nocase --names --file $GOFILE
package config

import "time"

type Config struct {
	RootPath  string `yaml:"root"` // server root folder
	Secret    string `yaml:"secret"`
	ServerURL string `yaml:"server_url"` // file browser url. required for ssl and oauth
	Locale    string `yaml:"locale"`     // default locale
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Auth      Auth   `yaml:"auth"`
	Store     Store  `yaml:"store"`
	SSL       SSL    `yaml:"ssl"`
}

type Auth struct {
	TTL struct {
		JWT    time.Duration `yaml:"jwt"`    // jwt TTL
		Cookie time.Duration `yaml:"cookie"` // auth cookie TTL
	} `yaml:"ttl"`
	Google    OAuth         `yaml:"google"`   // google oauth
	Github    OAuth         `yaml:"github"`   // github oauth
	Facebook  OAuth         `yaml:"facebook"` // facebook oauth
	Twitter   OAuth         `yaml:"twitter"`  // twitter oauth
	Dev       bool          `yaml:"dev"`      // enable dev (local) oauth2
	User      AuthUser      `yaml:"user"`
	Anonymous AnonymousUser `yaml:"anonymous"`
}

type OAuth struct {
	CID  string `yaml:"cid"`  // OAuth client ID
	CSEC string `yaml:"csec"` // OAuth client secret
}

type AuthUser struct {
	Home        string      `yaml:"home"`        // default user home
	Permissions Permissions `yaml:"permissions"` // default user permissions for home volume
}
type AnonymousUser struct {
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
statis
auto
)
*/
type SSLType int

type SSL struct {
	Type SSLType `yaml:"type"`
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

func defaultConfig() Config {
	return Config{
		RootPath: ".",
		Locale:   "en",
		Host:     "0.0.0.0",
		Port:     8080,
		Auth: Auth{
			TTL: struct {
				JWT    time.Duration `yaml:"jwt"`
				Cookie time.Duration `yaml:"cookie"`
			}{
				JWT:    10 * time.Second,
				Cookie: 7 * 24 * time.Hour,
			},
			User: AuthUser{
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
			Anonymous: AnonymousUser{
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
		SSL: SSL{
			Type: SSLTypeNone,
			Port: 8443,
		},
	}
}
