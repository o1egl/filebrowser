package config

import (
	"context"
	"io/ioutil"

	"github.com/goccy/go-yaml"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

type Loader interface {
	Load(ctx context.Context) ([]byte, error)
}

type LoaderFn func(ctx context.Context) ([]byte, error)

func (l LoaderFn) Load(ctx context.Context) ([]byte, error) {
	return l(ctx)
}

func FileLoader(filename string) Loader {
	return LoaderFn(func(_ context.Context) ([]byte, error) {
		return ioutil.ReadFile(filename)
	})
}

func Load(ctx context.Context, loader Loader) (*Config, error) {
	b, err := loader.Load(ctx)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, errors.Wrap(err, "failed to parse config file")
	}

	if err := mergo.Merge(&cfg, Default()); err != nil {
		return nil, err
	}

	if err := cfg.Validate(); err != nil {
		return nil, errors.Wrap(err, "config validation failed")
	}

	return &cfg, nil
}
