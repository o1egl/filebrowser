package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bep/simplecobra"
	"github.com/spf13/cobra"
)

var errHelp = fmt.Errorf("help")

func Execute(args []string) error {
	ctx, cancelCtx := context.WithCancel(context.Background())
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancelCtx()
	}()

	x, err := newExec()
	if err != nil {
		return err
	}
	cd, err := x.Execute(ctx, args)
	if err != nil {
		if errors.Is(err, errHelp) {
			_ = cd.CobraCommand.Help()
			fmt.Println()
			return nil
		}
		if simplecobra.IsCommandError(err) {
			// Print the help, but also return the error to fail the command.
			_ = cd.CobraCommand.Help()
			fmt.Println()
		}
	}
	return err
}

func newExec() (*simplecobra.Exec, error) {
	rootCmd := &rootCommand{
		commands: []simplecobra.Commander{
			newVersionCmd(),
			newServeCommand(),
		},
	}

	return simplecobra.New(rootCmd)
}

type rootCommand struct {
	commands []simplecobra.Commander
}

func (r *rootCommand) Name() string {
	return "filebrowser"
}

func (r *rootCommand) Init(commandeer *simplecobra.Commandeer) error {
	return nil
}

func (r *rootCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	return nil
}

func (r *rootCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	return errHelp
}

func (r *rootCommand) Commands() []simplecobra.Commander {
	return r.commands
}

type simpleCommand struct {
	use   string
	name  string
	short string
	long  string
	run   func(ctx context.Context, cd *simplecobra.Commandeer, rootCmd *rootCommand, args []string) error
	withc func(cmd *cobra.Command, r *rootCommand)
	initc func(cd *simplecobra.Commandeer) error

	commands []simplecobra.Commander

	rootCmd *rootCommand
}

func (c *simpleCommand) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *simpleCommand) Name() string {
	return c.name
}

func (c *simpleCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	if c.run == nil {
		return nil
	}
	return c.run(ctx, cd, c.rootCmd, args)
}

func (c *simpleCommand) Init(cd *simplecobra.Commandeer) error {
	c.rootCmd = cd.Root.Command.(*rootCommand)
	cmd := cd.CobraCommand
	cmd.Short = c.short
	cmd.Long = c.long
	if c.use != "" {
		cmd.Use = c.use
	}
	if c.withc != nil {
		c.withc(cmd, c.rootCmd)
	}
	return nil
}

func (c *simpleCommand) PreRun(cd, runner *simplecobra.Commandeer) error {
	if c.initc != nil {
		return c.initc(cd)
	}
	return nil
}
