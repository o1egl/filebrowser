//go:generate ${TOOLS_BIN}/go-enum --sql --marshal --nocase --names --file $GOFILE

package config

import (
	"github.com/filebrowser/filebrowser/errutils"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host  string `yaml:"host"`
	Port  int    `yaml:"port"`
	Log   Log    `yaml:"log"`
	Store Store  `yaml:"store"`
}

// ENUM(debug,info,warn,error)
type LogLevel string

// ENUM(text,json)
type LogFormat string

// ENUM(stdout,stderr)
type LogOutput string

type Log struct {
	Level  LogLevel  `yaml:"level"`
	Format LogFormat `yaml:"format"`
	Output LogOutput `yaml:"output"`
}

// ENUM(sqlite,mysql,postgres)
type StoreType int

type Store struct {
	Type     StoreType     `yaml:"type"`
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

// FromFile reads a config from a file.
func FromFile(fPath string) (_ *Config, err error) {
	file, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	defer errutils.JoinFn(&err, file.Close)

	cfg := Default()
	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Default returns a default config.
func Default() *Config {
	return &Config{
		Host: "localhost",
		Port: 8080,
		Log: Log{
			Level:  LogLevelInfo,
			Format: LogFormatText,
			Output: LogOutputStdout,
		},
		Store: Store{
			Type: StoreTypeSqlite,
			SQLite: SQLiteStore{
				File: "filebrowser.db",
			},
		},
	}
}
