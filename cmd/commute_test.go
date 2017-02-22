package cmd

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func TestCommuteCmd_Run(t *testing.T) {
	// Positive
	tests := []struct {
		from string
		to   string
	}{
		{from: "home", to: "work"},
		{from: "home", to: "123 Main st. Toronto Ontario"},
		{from: "321 Maple Ave. Toronto, Ontario", to: "123 Main st. Toronto Ontario"},
		{from: "321 Maple Ave. Toronto, Ontario", to: "work"},
	}

	for idx, tt := range tests {
		d := time.Minute * 3
		m := mockDurationer{
			durationFn: func(from, to string) (*time.Duration, error) {
				if from != tt.from {
					t.Fatalf("[#%v] Unexpected From, expected=%v, got=%v", idx, tt.from, from)
				} else if to != tt.to {
					t.Fatalf("[#%v] Unexpected To, expected=%v, got=%v", idx, tt.to, to)
				}

				return &d, nil
			},
		}

		c := CommuteCmd{From: tt.from, To: tt.to, Durationer: &m}
		var conf Configuration
		var i mockIndicator
		if err := c.Run(&conf, &i); err != nil {
			t.Fatal(err)
		}

		if len(i.out) != 1 {
			t.Fatalf("[#%v] Unexpected number of output lines, expected=%v, got=%v", idx, 1, i.out)
		} else if i.out[0] != c.format(d) {
			t.Fatalf("[#%v] Unexpected output, expected=%v, got=%v", idx, c.format(d), i.out[0])
		}
	}

	// Negative
	{
		testErr := errors.New("mock error")
		m := mockDurationer{
			durationFn: func(from, to string) (*time.Duration, error) {
				return nil, testErr
			},
		}

		c := CommuteCmd{From: "from", To: "to", Durationer: &m}
		var conf Configuration
		var i mockIndicator
		if err := c.Run(&conf, &i); err != testErr {
			t.Fatalf("Unexpected error returned, expected=%v, got=%v", testErr, err)
		}
	}
}

func TestCommuteCmd_format(t *testing.T) {
	tests := []struct {
		hours   int
		minutes int

		expected string
	}{
		{0, 0, "0 Minutes"},
		{1, 0, "1 Hour 0 Minutes"},
		{2, 1, "2 Hours 1 Minute"},
		{3, 3, "3 Hours 3 Minutes"},
	}

	var c CommuteCmd
	for _, tt := range tests {
		d := (time.Hour * time.Duration(tt.hours)) + (time.Minute * time.Duration(tt.minutes))

		out := c.format(d)
		if out != tt.expected {
			t.Fatalf("Unexpected output, expected=%v, got=%v", tt.expected, out)
		}
	}
}

func TestCommuteCmd_Validate(t *testing.T) {
	tests := []struct {
		locs        map[string]string
		from        string
		fromCurrent bool
		to          string
		toCurrent   bool
		err         error

		expectFrom string
		expectTo   string
	}{
		// Positive
		{map[string]string{}, "from", false, "to", false, nil, "from", "to"},
		{map[string]string{"from": "fromvalue"}, "from", false, "to", false, nil, "fromvalue", "to"},
		{map[string]string{"to": "tovalue"}, "from", false, "to", false, nil, "from", "tovalue"},
		{map[string]string{"from": "fromvalue", "to": "tovalue"}, "from", false, "to", false, nil, "fromvalue", "tovalue"},
		{map[string]string{}, "", true, "to", false, nil, "", "to"},
		{map[string]string{}, "from", false, "", true, nil, "from", ""},
		{map[string]string{}, "", true, "", true, nil, "", ""},

		// Negative
		{map[string]string{}, "", false, "to", false, ErrDefaultFromMissing, "", ""},
		{map[string]string{}, "from", false, "", false, ErrDefaultToMissing, "", ""},
		{map[string]string{}, "from", true, "", false, ErrFromAndFromCurrentProvided, "", ""},
		{map[string]string{}, "", false, "to", true, ErrToAndToCurrentProvided, "", ""},
		{map[string]string{"from": ""}, "from", false, "to", false, ErrDefaultFromMissing, "", ""},
		{map[string]string{"to": ""}, "from", false, "to", false, ErrDefaultToMissing, "", ""},
	}

	for idx, tt := range tests {
		conf := Configuration{Locations: tt.locs}
		m := mockLocator{
			locateFn: func() (float64, float64, error) {
				if !tt.fromCurrent && !tt.toCurrent {
					t.Fatalf("[#%v] Should not have called Locator", idx)
				}

				if tt.fromCurrent {
					tt.expectFrom = fmt.Sprintf("%v,%v", idx, idx)
				}
				if tt.toCurrent {
					tt.expectTo = fmt.Sprintf("%v,%v", idx, idx)
				}

				return float64(idx), float64(idx), nil
			},
		}
		c := CommuteCmd{From: tt.from, FromCurrent: tt.fromCurrent, To: tt.to, ToCurrent: tt.toCurrent, Locator: &m}

		if err := c.Validate(&conf); err != tt.err {
			t.Fatalf("[#%v] Unexpected error, expected=%v, got=%v", idx, tt.err, err)
		}

		if tt.err != nil {
			continue
		}

		if c.From != tt.expectFrom {
			t.Fatalf("[#%v] Unexpected From, expected=%v, got=%v", idx, tt.expectFrom, c.From)
		} else if c.To != tt.expectTo {
			t.Fatalf("[#%v] Unexpected To, expected=%v, got=%v", idx, tt.expectTo, c.To)
		}
	}
}
