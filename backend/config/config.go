//go:generate ${TOOLS_BIN}/go-enum --sql --marshal --nocase --names --file $GOFILE

package config

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/filebrowser/filebrowser/errutils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Host        string       `yaml:"host"`
	Port        int          `yaml:"port"`
	Domain      string       `yaml:"domain"`
	BasePath    string       `yaml:"basepath"`
	Log         Log          `yaml:"log"`
	Store       Store        `yaml:"store"`
	FileSystems []FileSystem `yaml:"filesystem"`
}

func (c Config) PublicAddress() string {
	address := net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
	if c.Domain != "" {
		address = c.Domain
	}
	return fmt.Sprintf("http://%s%s", address, c.BasePath)
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

// ENUM(local)
type FileSystemType int

type FileSystem struct {
	Name            string          `yaml:"name"`
	Type            FileSystemType  `yaml:"type"`
	LocalFileSystem LocalFileSystem `yaml:"local"`
}

type LocalFileSystem struct {
	Path string `yaml:"path"`
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
		Host:     "localhost",
		Port:     8080,
		BasePath: "/",
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
		FileSystems: []FileSystem{
			{
				Name: "local",
				Type: FileSystemTypeLocal,
				LocalFileSystem: LocalFileSystem{
					Path: "/",
				},
			},
		},
	}
}
