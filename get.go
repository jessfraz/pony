package main

import (
	"context"
	"errors"
	"flag"
	"fmt"

	"github.com/atotto/clipboard"
)

const getHelp = `Get details for a secret.`

func (cmd *getCommand) Name() string      { return "get" }
func (cmd *getCommand) Args() string      { return "[OPTIONS] KEY" }
func (cmd *getCommand) ShortHelp() string { return getHelp }
func (cmd *getCommand) LongHelp() string  { return getHelp }
func (cmd *getCommand) Hidden() bool      { return false }

func (cmd *getCommand) Register(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.copy, "copy", false, "copy the value to clipboard")
}

type getCommand struct {
	copy bool
}

func (cmd *getCommand) Run(ctx context.Context, args []string) error {
	if len(args) < 1 {
		return errors.New("must pass a key")
	}

	// Get the key value pair from secrets.
	key := args[0]
	if _, ok := s.Secrets[key]; !ok {
		return fmt.Errorf("secret for key %s does not exist", key)
	}

	fmt.Println(s.Secrets[key])

	if !cmd.copy {
		// Return early.
		return nil
	}

	// Copy to clipboard.
	if err := clipboard.WriteAll(s.Secrets[key]); err != nil {
		return fmt.Errorf("clipboard copy failed: %v", err)
	}
	fmt.Println("Copied to clipboard!")

	return nil
}
