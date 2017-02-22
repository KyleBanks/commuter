package cmd

import (
	"errors"
	"fmt"
)

var (
	// ErrAddNameMissing is returned when running the add command the the -name arugment is missing.
	ErrAddNameMissing = errors.New("missing -name parameter")
	// ErrAddLocationMissing is returned when running the add command the the -location arugment is missing.
	ErrAddLocationMissing = errors.New("missing -location parameter")
)

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
