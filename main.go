package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/genuinetools/pkg/cli"
	"github.com/jessfraz/pony/version"
	"github.com/sirupsen/logrus"
)

const (
	defaultFilestore string = ".pony"
	defaultGPGPath   string = ".gnupg/"
)

var (
	file    string
	gpgpath string
	keyid   string

	s secretFile

	debug bool
)

func main() {
	// Create a new cli program.
	p := cli.NewProgram()
	p.Name = "pony"
	p.Description = "Local File-Based Password, API Key, Secret, Recovery Code Store Backed By GPG"
	// Set the GitCommit and Version.
	p.GitCommit = version.GITCOMMIT
	p.Version = version.VERSION

	// Build the list of available commands.
	p.Commands = []cli.Command{
		&createCommand{},
		&getCommand{},
		&listCommand{},
		&removeCommand{},
	}

	// Setup the global flags.
	p.FlagSet = flag.NewFlagSet("global", flag.ExitOnError)
	p.FlagSet.StringVar(&file, "file", fmt.Sprintf("%s/%s", homeShortcut, defaultFilestore), "file to use for saving encrypted secrets")

	p.FlagSet.StringVar(&keyid, "keyid", os.Getenv("PONY_KEYID"), "optionally set specific gpg keyid/fingerprint to use for encryption & decryption (or env var PONY_KEYID)")

	p.FlagSet.BoolVar(&debug, "d", false, "enable debug logging")
	p.FlagSet.BoolVar(&debug, "debug", false, "enable debug logging")

	// Set the before function.
	p.Before = func(ctx context.Context) error {
		// Set the log level.
		if debug {
			logrus.SetLevel(logrus.DebugLevel)
		}

		home, err := getHome()
		if err != nil {
			logrus.Fatal(err)
		}

		// Set the file variable.
		file = strings.Replace(file, homeShortcut, home, 1)

		// Create our secrets file if it does not exist.
		if _, err := os.Stat(file); os.IsNotExist(err) {
			if err := writeSecretsFile(file, secretFile{}); err != nil {
				return err
			}
		}

		// We need to read the secrets file for all commands
		// might as well be dry about it.
		s, err = readSecretsFile(file)
		if err != nil {
			logrus.Fatal(err)
		}

		return nil
	}

	// Run our program.
	p.Run()
}

func getHome() (string, error) {
	home := os.Getenv(homeKey)
	if home != "" {
		return home, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.HomeDir, nil
}
