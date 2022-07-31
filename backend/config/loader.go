package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func LoadWithDefaults() (*Config, error) {
	cfg := &Config{
		Root:   ".",
		Secret: uuid.New().String(),
		Locale: "en",
		Upload: Upload{
			Path: "./var/uploads",
		},
		Server: Server{
			AccessLog: false,
			Bind:      ":8080",
		},
		Auth: Auth{
			TTL: struct {
				JWT    time.Duration `yaml:"jwt"`
				Cookie time.Duration `yaml:"cookie"`
			}{
				JWT:    10 * time.Second,
				Cookie: 7 * 24 * time.Hour,
			},
		},
		Store: Store{
			Type: StoreTypeSqlite,
			SQLite: SQLiteStore{
				File: "./var/filebrowser.db",
			},
		},
	}

	err := viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
