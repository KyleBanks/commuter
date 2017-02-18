package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/KyleBanks/commuter/pkg/storage"
)

const (
	msgGoogleMapsAPIKeyPrompt = "Enter Google Maps API Key: (developers.google.com/console)"
	msgDefaultLocationPrompt  = "Enter Your Default Location: (ex. 123 Main St. Toronto, Canada)"

	cmdDefault       string = "commuter"
	defaultFromParam string = "from"
	defaultFromUsage string = `The starting point of your commute, either a named location [ex. "work"] or an address [ex. "123 Main St. Toronto, Canada"].`
	defaultToParam   string = "to"
	defaultToUsage   string = `The destination of your commute, either a named location [ex. "work"] or an address [ex. "123 Main St. Toronto, Canada"].`

	cmdAdd           string = "add"
	addNameParam     string = "name"
	addNameUsage     string = `The name of the location you'd like to add [ex. "work"]. (required)`
	addLocationParam string = "location"
	addLocationUsage string = `The location to be added [ex. "123 Main St. Toronto, Canada"]. (required)`
)

// Stdout provides an output mechanism to notify the user via stdout.
type Stdout struct{}

// Indicate prints an indication to the user.
func (s Stdout) Indicate(v interface{}) {
	s.Indicatef("%v", v)
}

// Indicatef prints an indication to the user with formatting.
func (s Stdout) Indicatef(msg string, args ...interface{}) {
	fmt.Printf("> %v\n", fmt.Sprintf(msg, args...))
}

// Stdin provides an input mechanism for the user via the command line.
type Stdin struct {
	*bufio.Scanner
}

// newStdin initializes and returns a Stdin.
func newStdin() Stdin {
	return Stdin{
		Scanner: bufio.NewScanner(os.Stdin),
	}
}

// ArgParser parses input arguments from the command line.
type ArgParser struct {
	Args []string
}

// NewArgParser initializes and returns an ArgParser.
func NewArgParser(args []string) *ArgParser {
	return &ArgParser{
		Args: args,
	}
}

// Parse attempts to determine which command is being executed,
// parse its flags, and return it.
func (a *ArgParser) Parse(conf *Configuration, s storage.Provider) (Runner, error) {
	if conf == nil || len(a.Args) == 0 {
		return a.parseConfigureCmd(s)
	}

	switch a.Args[0] {
	case cmdAdd:
		return a.parseAddCmd(s, a.Args[1:])
	}

	return a.parseCommuteCmd(a.Args)
}

// parseConfigureCmd parses and returns a ConfigureCmd.
func (a *ArgParser) parseConfigureCmd(s storage.Provider) (*ConfigureCmd, error) {
	return &ConfigureCmd{
		Input: newStdin(),
		Store: s,
	}, nil
}

// parseCommuteCmd parses and returns a CommuteCmd from user supplied flags.
func (a *ArgParser) parseCommuteCmd(args []string) (*CommuteCmd, error) {
	var c CommuteCmd

	f := flag.NewFlagSet(cmdDefault, flag.ExitOnError)
	f.StringVar(&c.From, defaultFromParam, DefaultLocationAlias, defaultFromUsage)
	f.StringVar(&c.To, defaultToParam, DefaultLocationAlias, defaultToUsage)
	f.Parse(args)

	return &c, nil
}

// parseAddCmd parses and returns an AddCmd from user supplied flags.
func (a *ArgParser) parseAddCmd(s storage.Provider, args []string) (*AddCmd, error) {
	c := AddCmd{Store: s}

	f := flag.NewFlagSet(cmdAdd, flag.ExitOnError)
	f.StringVar(&c.Name, addNameParam, "", addNameUsage)
	f.StringVar(&c.Value, addLocationParam, "", addLocationUsage)
	f.Parse(args)

	return &c, nil
}
