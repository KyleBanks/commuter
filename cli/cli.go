// Package cli provides access to commuter over the command line via Stdin, Stdout,
// and command line arguments.
package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	cmdDefault       = "commuter"
	defaultFromParam = "from"
	defaultFromUsage = `The starting point of your commute, either a named location [ex. "work"] or an address [ex. "123 Main St. Toronto, Canada"].`
	defaultToParam   = "to"
	defaultToUsage   = `The destination of your commute, either a named location [ex. "work"] or an address [ex. "123 Main St. Toronto, Canada"].`

	cmdAdd           = "add"
	addNameParam     = "name"
	addNameUsage     = `The name of the location you'd like to add [ex. "work"]. (required)`
	addLocationParam = "location"
	addLocationUsage = `The location to be added [ex. "123 Main St. Toronto, Canada"]. (required)`
)

// Stdout provides an output mechanism to notify the user via stdout.
type Stdout struct {
	out io.Writer
}

// NewStdout initializes and returns a new Stdout.
func NewStdout() Stdout {
	return Stdout{
		out: os.Stdout,
	}
}

// Indicate prints an indication to the user.
func (s Stdout) Indicate(v interface{}) {
	s.Indicatef("%v", v)
}

// Indicatef prints an indication to the user with formatting.
func (s Stdout) Indicatef(msg string, args ...interface{}) {
	fmt.Fprintf(s.out, "%v\n", fmt.Sprintf(msg, args...))
}

// Stdin provides an input mechanism for the user via the command line.
type Stdin struct {
	*bufio.Scanner
}

// NewStdin initializes and returns a Stdin.
func NewStdin() Stdin {
	return Stdin{
		Scanner: bufio.NewScanner(os.Stdin),
	}
}
