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
	DefaultLocationAlias string = "default"

	// MsgGoogleMapsAPIKeyPrompt is used to prompt the user to enter their Google Maps API Key.
	MsgGoogleMapsAPIKeyPrompt = "Enter Google Maps API Key: (developers.google.com/console)"
	// MsgDefaultLocationPrompt is used to prompt the user to enter their default location.
	MsgDefaultLocationPrompt = "Enter Your Default Location: (ex. 123 Main St. Toronto, Canada)"
)

var (
	// ErrDefaultFromMissing is returned when running the default command the the -from arugment is missing.
	ErrDefaultFromMissing = errors.New("missing -from parameter")
	// ErrDefaultToMissing is returned when running the default command the the -to arugment is missing.
	ErrDefaultToMissing = errors.New("missing -to parameter")

	// ErrAddNameMissing is returned when running the add command the the -name arugment is missing.
	ErrAddNameMissing = errors.New("missing -name parameter")
	// ErrAddLocationMissing is returned when running the add command the the -location arugment is missing.
	ErrAddLocationMissing = errors.New("missing -location parameter")
)

// ConfigureCmd is used to configure the commuter application
type ConfigureCmd struct {
	Input Scanner
	Store StorageProvider
}

// Run prompts the user to configure the commuter application.
func (c *ConfigureCmd) Run(conf *Configuration, i Indicator) error {
	conf = &Configuration{
		APIKey: c.promptForString(i, MsgGoogleMapsAPIKeyPrompt),

		Locations: map[string]string{
			DefaultLocationAlias: c.promptForString(i, MsgDefaultLocationPrompt),
		},
	}

	return c.Store.Save(&conf)
}

// Validate validates the ConfigureCmd is properly initialized and ready to be Run.
func (c *ConfigureCmd) Validate(conf *Configuration) error {
	return nil
}

// promptForString prompts the user for a string input.
func (c *ConfigureCmd) promptForString(i Indicator, msg string) string {
	i.Indicate(msg)

	var in string
	for c.Input.Scan() {
		in = c.Input.Text()
		if len(in) == 0 {
			continue
		}

		break
	}

	return in
}

// String returns a string representation of the ConfigureCmd.
func (c *ConfigureCmd) String() string {
	return "Configure"
}

// CommuteCmd represents the standard command to
// retrieve the commute time between two locations.
type CommuteCmd struct {
	From string
	To   string
}

// Run calculates the distance between the From and To locations,
// and outputs the result.
func (c *CommuteCmd) Run(conf *Configuration, i Indicator) error {
	// TODO: This is too tightly coupled, it would be preferrable to
	// pass the router interface.
	r, err := geo.NewRouter(conf.APIKey)
	if err != nil {
		return err
	}

	d, err := r.Duration(c.From, c.To)
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
func (c *CommuteCmd) Validate(conf *Configuration) error {
	if val, ok := conf.Locations[c.From]; ok {
		c.From = val
	}
	if val, ok := conf.Locations[c.To]; ok {
		c.To = val
	}

	if len(c.From) == 0 {
		return ErrDefaultFromMissing
	}
	if len(c.To) == 0 {
		return ErrDefaultToMissing
	}

	return nil
}

// String returns a string representation of the CommuteCmd.
func (c *CommuteCmd) String() string {
	return fmt.Sprintf("From '%v' to '%v'", c.From, c.To)
}

// AddCmd represents a command to add a named location.
type AddCmd struct {
	Name  string
	Value string

	Store StorageProvider
}

// Run adds the named location, overwriting the existing value if necessary.
func (a *AddCmd) Run(conf *Configuration, i Indicator) error {
	conf.Locations[a.Name] = a.Value
	return a.Store.Save(conf)
}

// Validate validates the AddCmd is properly initialized and ready to be Run.
func (a *AddCmd) Validate(conf *Configuration) error {
	if len(a.Name) == 0 {
		return ErrAddNameMissing
	}
	if len(a.Value) == 0 {
		return ErrAddLocationMissing
	}

	return nil
}

// String returns a string representation of the AddCmd.
func (a *AddCmd) String() string {
	return fmt.Sprintf("Adding named location '%v' with value '%v'", a.Name, a.Value)
}
