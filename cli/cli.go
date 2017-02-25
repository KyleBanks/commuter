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
	cmdCommute              = "commuter"
	commuteFromParam        = "from"
	commuteFromUsage        = "The starting point of your commute, either a named location [ex. 'work'] or an address [ex. '123 Main St. Toronto, Canada']."
	commuteToParam          = "to"
	commuteToUsage          = "The destination of your commute, either a named location [ex. 'work'] or an address [ex. '123 Main St. Toronto, Canada']."
	commuteFromCurrentParam = "from-current"
	commuteFromCurrentUsage = "Sets your current location as the starting point of your commute. This uses Geolocation to attempt to determine your Latitude/Longitude based on IP Address. [Accuracy may vary]"
	commuteToCurrentParam   = "to-current"
	commuteToCurrentUsage   = "Sets your current location as the destination of your commute. This uses Geolocation to attempt to determine your Latitude/Longitude based on IP Address. [Accuracy may vary]"
	commuteDriveParam       = "drive"
	commuteDriveUsage       = "Adds 'driving' as a transit type [default]"
	commuteWalkParam        = "walk"
	commuteWalkUsage        = "Adds 'walking' as a transit type"
	commuteBikeParam        = "bike"
	commuteBikeUsage        = "Adds 'biking' as a transit type"
	commuteTransitParam     = "transit"
	commuteTransitUsage     = "Adds 'transit' as a transit type"

	cmdAdd           = "add"
	addNameParam     = "name"
	addNameUsage     = "The name of the location you'd like to add [ex. 'work']. (required)\n"
	addLocationParam = "location"
	addLocationUsage = "The location to be added [ex. '123 Main St. Toronto, Canada']. (required)\n"

	cmdList = "list"
)

// Stdout provides an output mechanism to notify the user via stdout.
type Stdout struct {
	io.Writer
}

// NewStdout initializes and returns a new Stdout.
func NewStdout() Stdout {
	return Stdout{
		Writer: os.Stdout,
	}
}

// Indicate prints an indication to the user.
func (s Stdout) Indicate(msg string, args ...interface{}) {
	fmt.Fprintf(s, "%v\n", fmt.Sprintf(msg, args...))
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
