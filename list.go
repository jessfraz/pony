package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"text/tabwriter"
)

const listHelp = `List secrets.`

func (cmd *listCommand) Name() string      { return "ls" }
func (cmd *listCommand) Args() string      { return "" }
func (cmd *listCommand) ShortHelp() string { return listHelp }
func (cmd *listCommand) LongHelp() string  { return listHelp }
func (cmd *listCommand) Hidden() bool      { return false }

func (cmd *listCommand) Register(fs *flag.FlagSet) {
	fs.StringVar(&cmd.filter, "f", "", "filter secrets keys by a regular expression")
	fs.StringVar(&cmd.filter, "filter", "", "filter secrets keys by a regular expression")
}

type listCommand struct {
	filter string
}

func (cmd *listCommand) Run(ctx context.Context, args []string) error {
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 3, ' ', 0)

	// print header
	fmt.Fprintln(w, "KEY\tVALUE")

	// print the keys alphabetically
	printSorted := func(m map[string]string) {
		mk := make([]string, len(m))
		i := 0
		for k := range m {
			mk[i] = k
			i++
		}
		sort.Strings(mk)

		for _, key := range mk {
			if len(cmd.filter) > 0 {
				if ok, _ := regexp.MatchString(cmd.filter, key); !ok {
					continue
				}
			}
			fmt.Fprintf(w, "%s\t%s\n", key, m[key])
		}
	}

	printSorted(s.Secrets)

	w.Flush()
	return nil
}
