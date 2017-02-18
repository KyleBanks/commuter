package main

import (
	"os"

	"github.com/KyleBanks/commuter/pkg/storage"
)

const (
	configurationFileName string = "commuter.conf"
)

func main() {
	var out Stdout
	store := storage.NewFileStore(configurationFileName)
	conf := NewConfiguration(store)
	parser := NewArgParser(os.Args[1:])

	r, err := parser.Parse(conf, store)
	if err != nil {
		panic(err)
	}

	exec(out, conf, r)
}

// exec executes a Runner with the Indicator and Configuration provided.
func exec(i Indicator, c *Configuration, r Runner) {
	err := r.Run(c, i)

	if err != nil {
		i.Indicatef("Command Failed: %v", r)
		i.Indicatef("Error:\n%v", err)
	}
}
