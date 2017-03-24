package main

import (
	"os"

	"github.com/KyleBanks/commuter/cli"
	"github.com/KyleBanks/commuter/cmd"
	"github.com/KyleBanks/go-kit/storage"
)

const (
	configurationFileName string = "config.json"
	configurationDirName  string = "commuter"
)

func main() {
	out := cli.NewStdout()
	store := storage.NewFileStore(configurationDirName, configurationFileName)
	conf := cmd.NewConfiguration(store)
	parser := cli.NewArgParser(os.Args[1:])

	r, err := parser.Parse(conf, store)
	if err != nil {
		panic(err)
	}

	exec(out, conf, r)
}

// exec validates and executes a Runner with the Indicator and Configuration provided.
func exec(i cmd.Indicator, c *cmd.Configuration, r cmd.RunnerValidator) {
	if err := r.Validate(c); err != nil {
		i.Indicate("Invalid command: %v", r)
		i.Indicate("Error: %v", err)
	}

	if err := r.Run(c, i); err != nil {
		i.Indicate("Command Failed: %v", r)
		i.Indicate("Error: %v", err)
	}
}
