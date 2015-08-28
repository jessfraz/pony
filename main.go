package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/Sirupsen/logrus"
	"github.com/atotto/clipboard"
	"github.com/codegangsta/cli"
	"github.com/docker/docker/pkg/homedir"
	"github.com/docker/docker/pkg/term"
)

const (
	defaultFilestore string = ".pony"
	defaultGPGPath   string = ".gnupg/"

	VERSION = "v0.1.0"
	BANNER  = ` _ __   ___  _ __  _   _ 
| '_ \ / _ \| '_ \| | | |
| |_) | (_) | | | | |_| |
| .__/ \___/|_| |_|\__, |
|_|                |___/
`
)

var (
	defaultGPGKey string
	filestore     string
	gpgPath       string
	publicKeyring string
	secretKeyring string

	s SecretFile

	debug   bool
	version bool
)

// preload initializes any global options and configuration
// before the main or sub commands are run
func preload(c *cli.Context) (err error) {
	if c.GlobalBool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	defaultGPGKey = c.GlobalString("keyid")

	home := homedir.Get()
	homeShort := homedir.GetShortcutString()

	// set the filestore variable
	filestore = strings.Replace(c.GlobalString("file"), homeShort, home, 1)

	// set gpg path variables
	gpgPath = strings.Replace(c.GlobalString("gpgpath"), homeShort, home, 1)
	publicKeyring = filepath.Join(gpgPath, "pubring.gpg")
	secretKeyring = filepath.Join(gpgPath, "secring.gpg")

	// if they passed an arguement, run the prechecks
	// TODO(jfrazelle): This will run even if the command they issue
	// does not exist, which is kinda shitty
	if len(c.Args()) > 0 {
		preChecks()
	}

	// we need to read the secrets file for all commands
	// might as well be dry about it
	s, err = readSecretsFile(filestore)
	if err != nil {
		logrus.Fatal(err)
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
		cli.StringFlag{
			Name:   "keyid",
			Usage:  "optionally set specific gpg keyid/fingerprint to use for encryption & decryption",
			EnvVar: fmt.Sprintf("%s_KEYID", strings.ToUpper(app.Name)),
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"save"},
			Usage:   "Add a new secret",
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 2 {
					logrus.Errorf("You need to pass a key and value to the command. ex: %s %s com.example.apikey EUSJCLLAWE", app.Name, c.Command.Name)
					cli.ShowSubcommandHelp(c)
				}

				// add the key value pair to secrets
				key, value := args[0], args[1]
				s.setKeyValue(key, value, false)

				fmt.Printf("Added %s %s to secrets", key, value)
			},
		},
		{
			Name:    "delete",
			Aliases: []string{"rm"},
			Usage:   "Delete a secret",
			Action: func(c *cli.Context) {
				args := c.Args()
				if len(args) < 1 {
					cli.ShowSubcommandHelp(c)
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
				// TODO(jfrazelle): add some filtering
				_, stdout, _ := term.StdStreams()
				w := tabwriter.NewWriter(stdout, 20, 1, 3, ' ', 0)

				// print header
				fmt.Fprintln(w, "KEY\tVALUE")

				// print the keys alphabetically
				printSorted := func(m map[string]string) {
					mk := make([]string, len(m))
					i := 0
					for k, _ := range m {
						mk[i] = k
						i++
					}
					sort.Strings(mk)

					for _, key := range mk {
						fmt.Fprintf(w, "%s\t%s\n", key, m[key])
					}
				}

				printSorted(s.Secrets)

				w.Flush()
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

				// add the key value pair to secrets
				key, value := args[0], args[1]
				s.setKeyValue(key, value, true)

				fmt.Printf("Updated secret %s to %s", key, value)
			},
		},
	}
	app.Run(os.Args)
}
