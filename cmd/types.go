package cmd

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
	Indicate(interface{})
	Indicatef(string, ...interface{})
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
