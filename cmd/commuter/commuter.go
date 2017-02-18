package main

import (
	"flag"

	"github.com/KyleBanks/commuter/pkg/storage"
)

const (
	configurationFileName string = "commuter.conf"
)

func main() {
	var out Logger
	store := storage.FileStore{
		Filename: configurationFileName,
	}
	conf := loadConfig(store)

	var r Runner
	if conf == nil {
		r = parseConfigureCmd(store)
	} else {
		r = parseCmd()
	}

	exec(out, conf, r)
}

// loadConfig attempts to load commuter configuration from the
// provided storage Provider.
//
// If an error occurs, the configuration is deemed corrupt and
// nil is returned.
func loadConfig(s storage.Provider) *Configuration {
	var c Configuration
	err := s.Load(&c)
	if err != nil {
		return nil
	}

	return &c
}

// parseCmd attempts to determine which command is being executed,
// parse its flags, and return it.
func parseCmd() Runner {
	return parseCommuteCmd()
}

// parseConfigureCmd parses and returns a ConfigureCmd.
func parseConfigureCmd(s storage.Provider) *ConfigureCmd {
	return &ConfigureCmd{
		Input: newStdin(),
		Store: s,
	}
}

// parseCommuteCmd parses and returns a CommuteCmd.
func parseCommuteCmd() *CommuteCmd {
	var c CommuteCmd
	flag.StringVar(&c.From, "from", DefaultLocationAlias, "The starting point of your commute.")
	flag.StringVar(&c.To, "to", "", "The destination of your commute.")
	flag.Parse()

	return &c
}

// exec executes a Runner with the Indicator and Configuration provided.
func exec(i Indicator, c *Configuration, r Runner) {
	err := r.Run(c, i)

	if err != nil {
		i.Indicatef("Command Failed: %v", r)
		i.Indicatef("Error:\n%v", err)
	}
}
