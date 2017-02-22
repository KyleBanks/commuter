package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// DefaultLocationAlias is the name of the default alias
	// used for 'From' addresses when one is not provided.
	DefaultLocationAlias = "default"
)

var (
	// ErrDefaultFromMissing is returned when running the default command and the -from argument is missing.
	ErrDefaultFromMissing = errors.New("missing -from or -from-current parameter")
	// ErrDefaultToMissing is returned when running the default command and the -to argument is missing.
	ErrDefaultToMissing = errors.New("missing -to or -to-current parameter")

	// ErrFromAndFromCurrentProvided is returned when the -from and -from-current arguments are both supplied.
	ErrFromAndFromCurrentProvided = errors.New("cannot use -from and -from-current arguments")
	// ErrToAndToCurrentProvided is returned when the -to and -to-current arguments are both supplied.
	ErrToAndToCurrentProvided = errors.New("cannot use -to and -to-current arguments")
)

// CommuteCmd represents the standard command to
// retrieve the commute time between two locations.
type CommuteCmd struct {
	From        string
	FromCurrent bool

	To        string
	ToCurrent bool

	Durationer Durationer
	Locator    Locator
}

// Run calculates the distance between the From and To locations,
// and outputs the result.
func (c *CommuteCmd) Run(conf *Configuration, i Indicator) error {
	d, err := c.Durationer.Duration(c.From, c.To)
	if err != nil {
		return err
	}

	i.Indicate(c.format(*d))
	return nil
}

// format takes a duration and returns a formatter representation.
func (c *CommuteCmd) format(d time.Duration) string {
	var out []string
	hours := int(d.Hours())
	minutes := int(d.Minutes())

	if hours > 0 {
		minutes -= (hours * 60)

		var suffix string
		if hours != 1 {
			suffix = "s"
		}

		out = append(out, fmt.Sprintf("%v Hour%v", hours, suffix))
	}

	var suffix string
	if minutes != 1 {
		suffix = "s"
	}
	out = append(out, fmt.Sprintf("%v Minute%v", minutes, suffix))

	return strings.Join(out, " ")
}

// Validate validates the CommuteCmd is properly initialized and ready to be Run.
//
// TODO: This is an ugly function, thing of a way of cleaning it up.
func (c *CommuteCmd) Validate(conf *Configuration) error {
	if c.FromCurrent && len(c.From) > 0 && c.From != DefaultLocationAlias {
		return ErrFromAndFromCurrentProvided
	}
	if c.ToCurrent && len(c.To) > 0 && c.To != DefaultLocationAlias {
		return ErrToAndToCurrentProvided
	}

	var err error
	c.From, err = c.getLoc(c.FromCurrent, c.From, conf)
	if err != nil {
		return err
	}
	c.To, err = c.getLoc(c.ToCurrent, c.To, conf)
	if err != nil {
		return err
	}

	if len(c.From) == 0 {
		return ErrDefaultFromMissing
	} else if len(c.To) == 0 {
		return ErrDefaultToMissing
	}

	return nil
}

// getLoc uses geolocation to populate a location if the current flag is true, otherwise
// it checks for a named location matching the input string.
//
// TODO: This is an ugly function, thing of a way of cleaning it up.
func (c *CommuteCmd) getLoc(current bool, input string, conf *Configuration) (string, error) {
	if !current {
		val, ok := conf.Locations[input]
		if ok {
			return val, nil
		}

		return input, nil
	}

	lat, long, err := c.Locator.CurrentLocation()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%v,%v", lat, long), nil
}

// String returns a string representation of the CommuteCmd.
func (c *CommuteCmd) String() string {
	return fmt.Sprintf("From '%v' to '%v'", c.From, c.To)
}
