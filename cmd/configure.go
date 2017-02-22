package cmd

const (
	// promptPrefix is the prefix using when prompting a user for input.
	promptPrefix = "> "

	// MsgGoogleMapsAPIKeyPrompt is used to prompt the user to enter their Google Maps API Key.
	MsgGoogleMapsAPIKeyPrompt = promptPrefix + "Enter Google Maps API Key: (developers.google.com/console)"
	// MsgDefaultLocationPrompt is used to prompt the user to enter their default location.
	MsgDefaultLocationPrompt = promptPrefix + "Enter Your Default Location: (ex. 123 Main St. Toronto, Canada)"
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
