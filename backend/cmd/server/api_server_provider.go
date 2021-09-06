package server

import (
	"errors"
	"net/url"

	"github.com/filebrowser/filebrowser/v3/rest/api"
	"github.com/google/wire"
)

var ApiServerSet = wire.NewSet(api.NewServer, SSLConfigProvider, ApiServerOptionsProvider)

func SSLConfigProvider(srvCmd *ServerCommand) (config api.SSLConfig, err error) {
	switch srvCmd.SSL.Type {
	case "none":
		config.SSLMode = api.None
	case "static":
		if srvCmd.SSL.Cert == "" {
			return config, errors.New("path to cert.pem is required")
		}
		if srvCmd.SSL.Key == "" {
			return config, errors.New("path to key.pem is required")
		}
		config.SSLMode = api.Static
		config.Port = srvCmd.SSL.Port
		config.Cert = srvCmd.SSL.Cert
		config.Key = srvCmd.SSL.Key
	case "auto":
		config.SSLMode = api.Auto
		config.Port = srvCmd.SSL.Port
		config.ACMELocation = srvCmd.SSL.ACMELocation
		if srvCmd.SSL.ACMEEmail != "" {
			config.ACMEEmail = srvCmd.SSL.ACMEEmail
		} else if u, e := url.Parse(srvCmd.ServerURL); e == nil {
			config.ACMEEmail = "admin@" + u.Hostname()
		}
	}
	return config, err
}

func ApiServerOptionsProvider(srvCmd *ServerCommand, sslConfig api.SSLConfig) api.Options {
	return api.Options{
		Host:         srvCmd.Host,
		Port:         srvCmd.Port,
		ServerURL:    srvCmd.ServerURL,
		SharedSecret: srvCmd.Secret,
		Revision:     srvCmd.Revision,
		AccessLog:    srvCmd.AccessLog.Enable,
		Anonymous:    srvCmd.Auth.Anonymous.Enable,
		SSLConfig:    sslConfig,
	}
}
