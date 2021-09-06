package cmd

import (
	"os"

	"github.com/filebrowser/filebrowser/v3/log"
)

// CommonOptionsCommander extends flags.Commander with SetCommon
// All commands should implement this interfaces
type CommonOptionsCommander interface {
	SetCommon(commonOpts CommonOpts)
	Execute(args []string) error
	HandleDeprecatedFlags() []DeprecatedFlag
}

// CommonOpts sets externally from main, shared across all commands
type CommonOpts struct {
	Revision string
}

// DeprecatedFlag contains information about deprecated option
type DeprecatedFlag struct {
	Old           string
	New           string
	RemoveVersion string
}

// SetCommon satisfies CommonOptionsCommander interface and sets common option fields
// The method called by main for each command
func (c *CommonOpts) SetCommon(commonOpts CommonOpts) {
	c.Revision = commonOpts.Revision
}

// HandleDeprecatedFlags sets new flags from deprecated and returns their list
func (c *CommonOpts) HandleDeprecatedFlags() []DeprecatedFlag { return nil }

// ResetEnv clears sensitive env vars
func ResetEnv(envs ...string) {
	for _, env := range envs {
		if err := os.Unsetenv(env); err != nil {
			log.Warnf("can't unset env %s, %s", env, err)
		}
	}
}
