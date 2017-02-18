package main

import (
	"fmt"

	"github.com/KyleBanks/commuter/pkg/geo"
	"github.com/KyleBanks/commuter/pkg/storage"
)

const (
	// DefaultLocationAlias is the name of the default alias
	// used for 'From' addresses when one is not provided.
	DefaultLocationAlias string = "default"
)

// ConfigureCmd is used to configure the commuter application
type ConfigureCmd struct {
	Input Scanner
	Store storage.Provider
}

// Run prompts the user to configure the commuter application.
func (c *ConfigureCmd) Run(conf *Configuration, i Indicator) error {
	conf = &Configuration{
		APIKey: c.promptForString(i, msgGoogleMapsAPIKeyPrompt),

		Locations: map[string]string{
			DefaultLocationAlias: c.promptForString(i, msgDefaultLocationPrompt),
		},
	}

	return c.Store.Save(&conf)
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
	r, err := geo.NewRouter(conf.APIKey)
	if err != nil {
		return err
	}

	if val, ok := conf.Locations[c.From]; ok {
		c.From = val
	}
	if val, ok := conf.Locations[c.To]; ok {
		c.To = val
	}

	d, err := r.Duration(c.From, c.To)
	if err != nil {
		return err
	}

	i.Indicate(d)
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

	Store storage.Provider
}

// Run adds the named location, overwriting the existing value if necessary.
func (a *AddCmd) Run(conf *Configuration, i Indicator) error {
	conf.Locations[a.Name] = a.Value
	return a.Store.Save(conf)
}

// String returns a string representation of the AddCmd.
func (a *AddCmd) String() string {
	return fmt.Sprintf("Adding named location '%v' with value '%v'", a.Name, a.Value)
}
