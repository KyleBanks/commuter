package cmd

import (
	"sort"
)

// ListCmd represents a request to list all locations.
type ListCmd struct{}

// Run lists all named aliases and their value.
func (l *ListCmd) Run(conf *Configuration, i Indicator) error {
	names := make([]string, len(conf.Locations), len(conf.Locations))
	var idx int
	var maxLen int
	for k := range conf.Locations {
		names[idx] = k

		if len(k) > maxLen {
			maxLen = len(k)
		}

		idx++
	}

	sort.Sort(byNameDefaultFirst(names))

	for _, name := range names {
		i.Indicate("%*s: %v", maxLen, name, conf.Locations[name])
	}

	return nil
}

// Validate ensures the ListCmd is valid to be executed.
func (l *ListCmd) Validate(conf *Configuration) error {
	return nil
}

// String returns a string representation of the ListCmd.
func (l *ListCmd) String() string {
	return "List"
}

// byNameDefaultFirst is used to sort a slice of strings alphabetically, with the 'default' alias
// always being first.
type byNameDefaultFirst []string

func (d byNameDefaultFirst) Len() int {
	return len(d)
}

func (d byNameDefaultFirst) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d byNameDefaultFirst) Less(i, j int) bool {
	if d[i] == DefaultLocationAlias {
		return true
	} else if d[j] == DefaultLocationAlias {
		return false
	}

	return d[i] < d[j]
}
