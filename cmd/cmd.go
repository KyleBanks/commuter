// Package cmd contains the core logic of commuter.
package cmd

import (
	"time"
)

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
func NewConfiguration(s StorageProvider) *Configuration {
	var c Configuration
	err := s.Load(&c)
	if err != nil {
		return nil
	}

	return &c
}

// Runner defines a type that can be Run.
type Runner interface {
	Run(*Configuration, Indicator) error
}

// Validator defines a type that can be validated.
type Validator interface {
	Validate(*Configuration) error
}

// RunnerValidator defines a type that can be Run and Validated.
type RunnerValidator interface {
	Runner
	Validator
}

// Indicator defines a type that can provide UI indications
// to the user.
type Indicator interface {
	Indicate(string, ...interface{})
}

// Scanner represents an input scanner that can read lines
// of input from the user.
type Scanner interface {
	Scan() bool
	Text() string
}

// StorageProvider defines a type that can be used for storage.
type StorageProvider interface {
	Load(interface{}) error
	Save(interface{}) error
}

// Durationer provides the ability to retrieve the duration between
// two locations.
type Durationer interface {
	Duration(string, string) (*time.Duration, error)
}

// Locator provides the ability to retrieve the current location as
// a Latitude and Longitude.
type Locator interface {
	CurrentLocation() (float64, float64, error)
}
