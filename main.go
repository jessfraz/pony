package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/homedir"
)

const (
	defaultFilestore string = ".pony"
	defaultGPGPath   string = ".gnupg"

	VERSION = "v0.1.0"
	BANNER  = ` _ __   ___  _ __  _   _ 
| '_ \ / _ \| '_ \| | | |
| |_) | (_) | | | | |_| |
| .__/ \___/|_| |_|\__, |
|_|                |___/
`
)

var (
	filestore     string
	gpgPath       string
	publicKeyring string
	secretKeyring string

	debug   bool
	version bool
)

// preload initializes any global options and configuration
// before the main or sub commands are run
func preload(c *cli.Context) error {
	if c.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// set the filestore variable
	filestore = filepath.Join(homedir.Get(), strings.TrimPrefix(c.GlobalString("file"), homedir.GetShortcutString()))

	// set gpg path variables
	gpgPath = filepath.Join(homedir.Get(), strings.TrimPrefix(c.GlobalString("gpgpath"), homedir.GetShortcutString()))
	publicKeyring = filepath.Join(gpgPath, "pubring.gpg")
	secretKeyring = filepath.Join(gpgPath, "secring.gpg")

	// if they passed an arguement, run the prechecks
	// TODO(jfrazelle): This will run even if the command they issue
	// does not exist, which is kinda shitty
	if len(c.Args()) > 0 {
		preChecks()
	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "pony"
	app.Version = VERSION
	app.Author = "@jfrazelle"
	app.Email = "no-reply@butts.com"
	app.Usage = "Local File-Based Password, API Key, Secret, Recovery Code Store Backed By GPG"
	app.Before = preload
	app.EnableBashCompletion = true
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "run in debug mode",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: fmt.Sprintf("%s/%s", homedir.GetShortcutString(), defaultFilestore),
			Usage: "file to use for saving encrypted secrets",
		},
		cli.StringFlag{
			Name:  "gpgpath",
			Value: fmt.Sprintf("%s/%s", homedir.GetShortcutString(), defaultGPGPath),
			Usage: "filepath used for gpg keys",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"save"},
			Usage:   "Add a new secret",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type",
					Value: "kv",
					Usage: "Type of secret, key value (kv) or recovery code store (rcs)",
				},
			},
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 2 {
					logrus.Errorf("You need to pass a key and value to the command. ex: %s %s com.example.apikey EUSJCLLAWE", app.Name, c.Command.Name)
					cli.ShowSubcommandHelp(c)
				}

				s, err := readSecretsFile(filestore)
				if err != nil {
					logrus.Fatal(err)
				}

				// add the key value pair to secrets
				key, value := args[0], args[1]
				s.setKeyValue(key, value, false)

				fmt.Printf("Added %s %s to secrets\n", key, value)
			},
		},
		{
			Name:  "delete",
			Usage: "Delete a secret or a specific recovery code",
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 1 {
					cli.ShowSubcommandHelp(c)
				}

				s, err := readSecretsFile(filestore)
				if err != nil {
					logrus.Fatal(err)
				}

				key := args[0]
				if _, ok := s.Secrets[key]; !ok {
					logrus.Fatalf("Secret for (%s) does not exist", key)
				}
				delete(s.Secrets, key)

				if err := writeSecretsFile(filestore, s); err != nil {
					logrus.Fatal(err)
				}

				fmt.Printf("Secret %q deleted successfully", key)
			},
		},
		{
			Name:  "get",
			Usage: "Get the value of a secret",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "copy, c",
					Usage: "copy the secret to your clipboard",
				},
			},
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 1 {
					cli.ShowSubcommandHelp(c)
				}

				s, err := readSecretsFile(filestore)
				if err != nil {
					logrus.Fatal(err)
				}

				// add the key value pair to secrets
				key := args[0]
				if _, ok := s.Secrets[key]; !ok {
					logrus.Fatalf("Secret for (%s) does not exist", key)
				}

				fmt.Println(s.Secrets[key])

				// copy to clipboard
				if c.Bool("copy") {
					if err := clipboard.WriteAll(s.Secrets[key]); err != nil {
						logrus.Fatal("Clipboard copy failed: %v", err)
					}
					fmt.Println("Copied to clipboard!")
				}
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List all secrets",
			Action: func(c *cli.Context) {
				s, err := readSecretsFile(filestore)
				if err != nil {
					logrus.Fatal(err)
				}

				fmt.Printf("secrets file: %+v", s)
			},
		},
		{
			Name:  "update",
			Usage: "Update a secret",
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 2 {
					logrus.Errorf("You need to pass a key and value to the command. ex: %s %s com.example.apikey EUSJCLLAWE", app.Name, c.Command.Name)
					cli.ShowSubcommandHelp(c)
				}

				s, err := readSecretsFile(filestore)
				if err != nil {
					logrus.Fatal(err)
				}

				// add the key value pair to secrets
				key, value := args[0], args[1]
				s.setKeyValue(key, value, true)

				fmt.Printf("Updated secret %s to %s\n", key, value)
			},
		},
	}
	app.Run(os.Args)
}
