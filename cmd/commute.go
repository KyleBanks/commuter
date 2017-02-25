package cmd

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KyleBanks/commuter/pkg/geo"
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
	// ErrNoCommuteMethod is returned when no commute method is selected.
	ErrNoCommuteMethod = errors.New("at least one commute method must be specified")

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

	Drive   bool
	Walk    bool
	Bike    bool
	Transit bool

	Durationer Durationer
	Locator    Locator
}

// Run calculates the distance between the From and To locations,
// and outputs the result.
func (c *CommuteCmd) Run(conf *Configuration, i Indicator) error {
	modes := c.modes()
	multiMode := len(modes) > 1

	for _, m := range modes {
		d, err := c.Durationer.Duration(c.From, c.To, m)
		if err != nil {
			return err
		}

		var method string
		if multiMode {
			method = fmt.Sprintf("%v: ", m)
		}
		i.Indicate("%v%v", method, c.format(*d))
	}

	return nil
}

// format takes a duration and returns a formatter representation.
func (c *CommuteCmd) format(d time.Duration) string {
	pluralize := func(s string, i int) string {
		if i != 1 {
			s += "s"
		}
		return s
	}

	var out []string
	hours := int(d.Hours())
	minutes := int(d.Minutes())

	if hours > 0 {
		minutes -= (hours * 60)
		out = append(out, fmt.Sprintf("%v %v", hours, pluralize("Hour", hours)))
	}

	out = append(out, fmt.Sprintf("%v %v", minutes, pluralize("Minute", minutes)))

	return strings.Join(out, " ")
}

// modes returns a slice of TravelModes based on the commands Drive, Walk, Bike and Transit properties.
func (c *CommuteCmd) modes() []geo.TravelMode {
	var modes []geo.TravelMode
	if c.Drive {
		modes = append(modes, geo.Drive)
	}
	if c.Walk {
		modes = append(modes, geo.Walk)
	}
	if c.Bike {
		modes = append(modes, geo.Bike)
	}
	if c.Transit {
		modes = append(modes, geo.Transit)
	}

	return modes
}

// Validate validates the CommuteCmd is properly initialized and ready to be Run.
func (c *CommuteCmd) Validate(conf *Configuration) (err error) {
	// Must provide at least one method of transport
	if !c.Drive && !c.Walk && !c.Bike && !c.Transit {
		return ErrNoCommuteMethod
	}

	c.From, err = c.setLocation(conf, c.From, c.FromCurrent, ErrFromAndFromCurrentProvided, ErrDefaultFromMissing)
	if err != nil {
		return
	}

	c.To, err = c.setLocation(conf, c.To, c.ToCurrent, ErrToAndToCurrentProvided, ErrDefaultToMissing)
	if err != nil {
		return
	}

	return
}

// setLocation validates and determines a location based on the provided value and the `useCurrent` flag.
//
// If the useCurrent flag is true, setLocation will attempt to use geolocation to determine the current location. Otherwise,
// it will check if the value is an alias, or use the actual value provided.
func (c *CommuteCmd) setLocation(conf *Configuration, value string, useCurrent bool, bothProvided error, missing error) (string, error) {
	if useCurrent && len(value) > 0 && value != DefaultLocationAlias {
		return "", bothProvided
	}

	var err error
	if useCurrent {
		value, err = c.locate(conf)
	} else {
		value = c.alias(conf, value)
	}

	if err != nil {
		return "", err
	}

	if len(value) == 0 {
		return "", missing
	}

	return value, nil
}

// alias checks if the provided value is an alias to a location in the Configuration.
// If it is, the value of the alias is returned, otherwise the value provided is returned.
func (c *CommuteCmd) alias(conf *Configuration, value string) string {
	val, ok := conf.Locations[value]
	if !ok {
		return value
	}

	return val
}

// locate attempts to return a latitude/longitude string for the user's current location.
func (c *CommuteCmd) locate(conf *Configuration) (string, error) {
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
