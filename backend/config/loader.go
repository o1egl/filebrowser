package config

import (
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func Load() (*Config, error) {
	var cfg Config
	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadWithDefaults() (*Config, error) {
	viper.SetDefault("secret", uuid.New().String())
	viper.SetDefault("root", ".")
	viper.SetDefault("locale", "en")
	viper.SetDefault("server",
		Server{
			AccessLog: false,
			Bind:      ":8080",
		},
	)
	viper.SetDefault("upload",
		Upload{
			Path: "./var/uploads",
		},
	)
	viper.SetDefault("store",
		Store{
			Type: StoreTypeSqlite,
			SQLite: SQLiteStore{
				File: "./var/filebrowser.db",
			},
		},
	)
	viper.SetDefault("auth",
		Auth{
			TTL: struct {
				JWT    time.Duration `yaml:"jwt"`
				Cookie time.Duration `yaml:"cookie"`
			}{
				JWT:    10 * time.Second,
				Cookie: 7 * 24 * time.Hour,
			},
			User: User{
				GenerateHome: true,
				Home: HomeVolume{
					Path: "/user",
					Permissions: VolumePermissions{
						Read:   true,
						Create: true,
						Modify: true,
						Delete: true,
						Share:  true,
					},
				},
			},
			Anonymous: Anonymous{
				Enabled: false,
				Home: HomeVolume{
					Path: "/",
					Permissions: VolumePermissions{
						Read:   true,
						Create: false,
						Modify: false,
						Delete: false,
						Share:  false,
					},
				},
			},
		},
	)
	return Load()
}
