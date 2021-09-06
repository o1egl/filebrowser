package server

import (
	"fmt"
	"net/url"
	"os"

	"github.com/filebrowser/filebrowser/v3/rest/api"
	"github.com/pkg/errors"
)

// serverApp holds all active objects
type serverApp struct {
	*ServerCommand
	restSrv    *api.Server
	terminated chan struct{}
}

func NewServerApp(srvCmd *ServerCommand, restSrv *api.Server) (*serverApp, error) {
	return &serverApp{
		ServerCommand: srvCmd,
		restSrv:       restSrv,
		terminated:    make(chan struct{}),
	}, nil
}

func NewSSLConfig(srvCmd *ServerCommand) (config api.SSLConfig, err error) {
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

// mkdir -p for all dirs
func makeDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil { // If path is already a directory, MkdirAll does nothing
			return fmt.Errorf("can't make directory %s: %w", dir, err)
		}
	}
	return nil
}
