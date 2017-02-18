package main

import (
	"os"

	"github.com/KyleBanks/commuter/cli"
	"github.com/KyleBanks/commuter/cmd"
	"github.com/KyleBanks/commuter/pkg/storage"
)

const (
	configurationFileName string = "commuter.conf"
)

func main() {
	out := cli.NewStdout()
	store := storage.NewFileStore(configurationFileName)
	conf := cmd.NewConfiguration(store)
	parser := cli.NewArgParser(os.Args[1:])

	r, err := parser.Parse(conf, store)
	if err != nil {
		panic(err)
	}

	exec(out, conf, r)
}

// exec validates executes a Runner with the Indicator and Configuration provided.
func exec(i cmd.Indicator, c *cmd.Configuration, r cmd.RunnerValidator) {
	if err := r.Validate(c); err != nil {
		i.Indicatef("Invalid command: %v", r)
		i.Indicatef("Error: %v", err)
	}

	if err := r.Run(c, i); err != nil {
		i.Indicatef("Command Failed: %v", r)
		i.Indicatef("Error: %v", err)
	}
}
