package main

import (
	"github.com/KyleBanks/commuter/pkg/storage"
)

// Runner defines a type that can be Run.
type Runner interface {
	Run(*Configuration, Indicator) error
}

// Indicator defines a type that can provide UI indications
// to the user.
type Indicator interface {
	Indicate(interface{})
	Indicatef(string, ...interface{})
}

// Scanner represents an input scanner that can read lines
// of input from the user.
type Scanner interface {
	Scan() bool
	Text() string
}

// Configuration represents a Commuter configuration, including
// the Google Maps API Key and location map.
type Configuration struct {
	APIKey    string
	Locations map[string]string
}

// NewConfiguration attempts to retrieve a Configuration from a storage Provider.
//
// If an error occurs, the configuration is assumed to be corrupted and a nil
// Configuration is returned.
func NewConfiguration(s storage.Provider) *Configuration {
	var c Configuration
	err := s.Load(&c)
	if err != nil {
		return nil
	}

	return &c
}
