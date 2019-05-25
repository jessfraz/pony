package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
)

const removeHelp = `Delete a secret.`

func (cmd *removeCommand) Name() string      { return "rm" }
func (cmd *removeCommand) Args() string      { return "[OPTIONS] KEY" }
func (cmd *removeCommand) ShortHelp() string { return removeHelp }
func (cmd *removeCommand) LongHelp() string  { return removeHelp }
func (cmd *removeCommand) Hidden() bool      { return false }

func (cmd *removeCommand) Register(fs *flag.FlagSet) {}

type removeCommand struct{}

func (cmd *removeCommand) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("must pass a key")
	}

	key := args[0]
	if _, ok := s.Secrets[key]; !ok {
		return fmt.Errorf("secret for key %s does not exist", key)
	}
	delete(s.Secrets, key)

	if err := writeSecretsFile(file, s); err != nil {
		return err
	}

	fmt.Printf("Deleted secret key %s\n", key)
	return nil
}
