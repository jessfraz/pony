package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

const createHelp = `Create a secret.`

func (cmd *createCommand) Name() string      { return "create" }
func (cmd *createCommand) Args() string      { return "[OPTIONS] KEY VALUE" }
func (cmd *createCommand) ShortHelp() string { return createHelp }
func (cmd *createCommand) LongHelp() string  { return createHelp }
func (cmd *createCommand) Hidden() bool      { return false }

func (cmd *createCommand) Register(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.force, "force", false, "force overwrite existing value")
	fs.BoolVar(&cmd.force, "f", false, "force overwrite existing value")
}

type createCommand struct {
	force bool
}

func (cmd *createCommand) Run(ctx context.Context, args []string) error {
	if len(args) < 2 {
		return errors.New("must pass a key and value")
	}

	// Check if we are updating.
	verb := "Added"
	_, isUpdating := s.Secrets[args[0]]
	if isUpdating {
		verb = "Updated"
	}

	// Add the key value pair to secrets.
	key, value := args[0], args[1]
	s.setKeyValue(key, value, cmd.force)

	fmt.Printf("%s %s %s to secrets\n", verb, key, value)
	return nil
}
