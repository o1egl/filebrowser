//go:generate ${TOOLS_BIN}/go-enum --sql --marshal --nocase --names --file $GOFILE

package config

import (
	"github.com/filebrowser/filebrowser/errutils"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Log  Log    `yaml:"log"`
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
	}
}
