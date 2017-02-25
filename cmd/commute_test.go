package cmd

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/KyleBanks/commuter/pkg/geo"
)

func TestCommuteCmd_Run(t *testing.T) {
	// Positive, To/From
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
			durationFn: func(from, to string, tm geo.TravelMode) (*time.Duration, error) {
				if from != tt.from {
					t.Fatalf("[#%v] Unexpected From, expected=%v, got=%v", idx, tt.from, from)
				} else if to != tt.to {
					t.Fatalf("[#%v] Unexpected To, expected=%v, got=%v", idx, tt.to, to)
				}

				return &d, nil
			},
		}

		c := CommuteCmd{From: tt.from, To: tt.to, Drive: true, Durationer: &m}
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

	// Positive, Modes
	mTests := []struct {
		drive   bool
		walk    bool
		bike    bool
		transit bool
	}{
		{true, false, false, false},
		{true, true, false, false},
		{false, true, false, false},
		{false, true, true, false},
		{false, false, true, false},
		{false, false, true, true},
		{false, true, true, true},
		{true, true, true, true},
	}

	for idx, tt := range mTests {
		var expectCount int
		if tt.drive {
			expectCount++
		}
		if tt.walk {
			expectCount++
		}
		if tt.bike {
			expectCount++
		}
		if tt.transit {
			expectCount++
		}

		m := mockDurationer{
			durationFn: func(from, to string, tm geo.TravelMode) (*time.Duration, error) {
				if tm == geo.Drive && !tt.drive {
					t.Fatalf("[#%v] Unexpected Mode, Drive", idx)
				} else if tm == geo.Walk && !tt.walk {
					t.Fatalf("[#%v] Unexpected Mode, Walk", idx)
				} else if tm == geo.Bike && !tt.bike {
					t.Fatalf("[#%v] Unexpected Mode, Bike", idx)
				} else if tm == geo.Transit && !tt.transit {
					t.Fatalf("[#%v] Unexpected Mode, Transit", idx)
				}

				d := time.Minute * 3
				return &d, nil
			},
		}

		c := CommuteCmd{From: "default", To: "default", Drive: tt.drive, Walk: tt.walk, Bike: tt.bike, Transit: tt.transit, Durationer: &m}
		var conf Configuration
		var i mockIndicator
		if err := c.Run(&conf, &i); err != nil {
			t.Fatal(err)
		}

		if len(i.out) != expectCount {
			t.Fatalf("[#%v] Unexpected number of output lines, expected=%v, got=%v", idx, expectCount, i.out)
		}
	}

	// Negative
	{
		testErr := errors.New("mock error")
		m := mockDurationer{
			durationFn: func(from, to string, tm geo.TravelMode) (*time.Duration, error) {
				return nil, testErr
			},
		}

		c := CommuteCmd{From: "from", To: "to", Drive: true, Durationer: &m}
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
	// From/To
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
		{map[string]string{}, "from", false, "to", true, ErrToAndToCurrentProvided, "", ""},
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
		c := CommuteCmd{From: tt.from, FromCurrent: tt.fromCurrent, To: tt.to, ToCurrent: tt.toCurrent, Locator: &m, Drive: true}

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

	// Modes
	cTests := []struct {
		drive   bool
		walk    bool
		bike    bool
		transit bool

		err error
	}{
		{true, false, false, false, nil},
		{true, true, false, false, nil},
		{true, true, true, false, nil},
		{true, true, true, true, nil},
		{false, true, true, true, nil},
		{false, false, true, true, nil},
		{false, false, false, true, nil},
		{false, false, false, false, ErrNoCommuteMethod},
	}

	for idx, tt := range cTests {
		conf := Configuration{Locations: make(map[string]string)}

		c := CommuteCmd{From: "default", To: "default", Drive: tt.drive, Walk: tt.walk, Bike: tt.bike, Transit: tt.transit}
		if err := c.Validate(&conf); err != tt.err {
			t.Fatalf("[#%v] Unexpected error, expected=%v, got=%v", idx, tt.err, err)
		}
	}
}
