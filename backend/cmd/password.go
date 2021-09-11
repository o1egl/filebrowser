package cmd

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/filebrowser/filebrowser/v3/hash"
)

// PasswordCommand with command line flags and env
type PasswordCommand struct {
	Secret string `long:"secret" env:"SECRET" description:"shared secret key"`
	CommonOpts
}

// Execute runs file browser server
func (s *PasswordCommand) Execute(args []string) error {
	switch {
	case len(args) == 0:
		return errors.New("password is not provided")
	case len(args) > 1:
		return errors.New("more than 1 argument provided")
	}

	password, err := hash.NewHasher(s.Secret).Password(args[0])
	if err != nil {
		return err
	}
	fmt.Println(password)

	return nil
}
