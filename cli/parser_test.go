package cli

import (
	"reflect"
	"testing"

	"github.com/KyleBanks/commuter/cmd"
)

type MockStorageProvider struct {
	loadFn func(interface{}) error
	saveFn func(interface{}) error
}

func (m MockStorageProvider) Load(i interface{}) error {
	return m.loadFn(i)
}
func (m MockStorageProvider) Save(i interface{}) error {
	return m.saveFn(i)
}

func TestNewArgParser(t *testing.T) {
	tests := []struct {
		args []string
	}{
		{nil},
		{[]string{}},
		{[]string{"Arg"}},
		{[]string{"Arg1", "Arg2"}},
	}

	for _, tt := range tests {
		p := NewArgParser(tt.args)
		if !testStringsEq(p.Args, tt.args) {
			t.Fatalf("Unexpected Args, expected=%v, got=%v", tt.args, p.Args)
		}
	}
}

func TestArgParser_Parse(t *testing.T) {
	conf := cmd.Configuration{APIKey: "123"}
	var s MockStorageProvider

	// Valid state and args
	tests := []struct {
		args     []string
		conf     *cmd.Configuration
		expected interface{}
	}{
		// Default command
		{[]string{"-to", "work"}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-from", "work"}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-to", "work", "-from", "home"}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-to", "123 Etc Ave.", "-from", "home"}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-to", "work", "-from", "123 Etc Ave."}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-to", "321 Example Drive", "-from", "123 Etc Ave."}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-to", "321 Example Drive", "-from-current"}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-to-current", "-from", "123 Etc Ave."}, &conf, &cmd.CommuteCmd{}},
		{[]string{"-to-current", "-from-current"}, &conf, &cmd.CommuteCmd{}},

		// Add command
		{[]string{"add"}, &conf, &cmd.AddCmd{}},
		{[]string{"add", "-name", "work"}, &conf, &cmd.AddCmd{}},
		{[]string{"add", "-name", "work", "-location", "123 Sample Lane"}, &conf, &cmd.AddCmd{}},

		// List command
		{[]string{"list"}, &conf, &cmd.ListCmd{}},
		{[]string{"list", "-arg"}, &conf, &cmd.ListCmd{}},

		// Empty args should prompt a ConfigureCommand
		{[]string{}, &conf, &cmd.ConfigureCmd{}},

		// Nil configuration should always prompt a ConfigureCmd
		{[]string{}, nil, &cmd.ConfigureCmd{}},
		{[]string{"-help"}, nil, &cmd.ConfigureCmd{}},
		{[]string{"-to", "work"}, nil, &cmd.ConfigureCmd{}},
		{[]string{"add"}, nil, &cmd.ConfigureCmd{}},
	}

	for idx, tt := range tests {
		p := NewArgParser(tt.args)
		r, err := p.Parse(tt.conf, s)
		if err != nil {
			t.Fatal(err)
		}

		if reflect.TypeOf(r) != reflect.TypeOf(tt.expected) {
			t.Fatalf("[%v] Unexpected Command parsed, expected=%v, got=%v", idx, reflect.TypeOf(tt.expected), reflect.TypeOf(r))
		}
	}
}

func TestArgParser_parseConfigureCmd(t *testing.T) {
	var s MockStorageProvider
	var a ArgParser

	c, err := a.parseConfigureCmd(&s)
	if err != nil {
		t.Fatal(err)
	}

	if c.Input == nil {
		t.Fatalf("Unexpected nil input")
	}
	if c.Store != &s {
		t.Fatalf("Unexpected Store, expected=%v, got=%v", s, c.Store)
	}
}

func TestArgParser_parseCommuteCmd(t *testing.T) {
	var conf cmd.Configuration
	var a ArgParser

	// No API key should return an error
	_, err := a.parseCommuteCmd(&conf, []string{"-to", "work"})
	if err == nil {
		t.Fatalf("Expected error for empty API key")
	}

	conf.APIKey = "example"

	// To/From and ToCurrent/FromCurrent
	tests := []struct {
		args     []string
		expected cmd.CommuteCmd
	}{
		{[]string{"-to", "work"}, cmd.CommuteCmd{To: "work", From: "default"}},
		{[]string{"-from", "work"}, cmd.CommuteCmd{To: "default", From: "work"}},
		{[]string{"-to", "work", "-from", "home"}, cmd.CommuteCmd{To: "work", From: "home"}},
		{[]string{"-to", "123 Etc Ave.", "-from", "home"}, cmd.CommuteCmd{To: "123 Etc Ave.", From: "home"}},
		{[]string{"-to", "work", "-from", "123 Etc Ave."}, cmd.CommuteCmd{To: "work", From: "123 Etc Ave."}},
		{[]string{"-to", "321 Example Drive", "-from", "123 Etc Ave."}, cmd.CommuteCmd{To: "321 Example Drive", From: "123 Etc Ave."}},
		{[]string{"-to", "321 Example Drive", "-from-current"}, cmd.CommuteCmd{To: "321 Example Drive", FromCurrent: true, From: "default"}},
		{[]string{"-to-current", "-from", "123 Etc Ave."}, cmd.CommuteCmd{ToCurrent: true, From: "123 Etc Ave.", To: "default"}},
		{[]string{"-to-current", "-from-current"}, cmd.CommuteCmd{ToCurrent: true, FromCurrent: true, To: "default", From: "default"}},
	}

	for idx, tt := range tests {
		r, err := a.parseCommuteCmd(&conf, tt.args)
		if err != nil {
			t.Fatal(err)
		}

		if tt.expected.To != r.To {
			t.Fatalf("[%v] Unexpected 'To' parsed, expected=%v, got=%v", idx, tt.expected.To, r.To)
		} else if tt.expected.From != r.From {
			t.Fatalf("[%v] Unexpected 'From' parsed, expected=%v, got=%v", idx, tt.expected.From, r.From)
		} else if tt.expected.FromCurrent != r.FromCurrent {
			t.Fatalf("[%v] Unexpected 'FromCurrent' parsed, expected=%v, got=%v", idx, tt.expected.FromCurrent, r.FromCurrent)
		} else if tt.expected.ToCurrent != r.ToCurrent {
			t.Fatalf("[%v] Unexpected 'ToCurrent' parsed, expected=%v, got=%v", idx, tt.expected.ToCurrent, r.ToCurrent)
		} else if r.Durationer == nil {
			t.Fatalf("[%v] Unexpected nil Durationer", idx)
		} else if r.Locator == nil {
			t.Fatalf("[%v] Unexpected nil Locator", idx)
		} else if r.Drive == false {
			t.Fatalf("[%v] Unexpected Drive, expected=true, got=false", idx)
		} else if r.Bike == true || r.Walk == true || r.Transit == true {
			t.Fatalf("[%v] Unexpected Bike/Walk/Transit, expected=false, got=[%v, %v, %v]", idx, r.Bike, r.Walk, r.Transit)
		}
	}

	// Modes
	mTests := []struct {
		args     []string
		expected cmd.CommuteCmd
	}{
		{[]string{""}, cmd.CommuteCmd{Drive: true}},
		{[]string{"-drive"}, cmd.CommuteCmd{Drive: true}},
		{[]string{"-walk"}, cmd.CommuteCmd{Walk: true}},
		{[]string{"-bike"}, cmd.CommuteCmd{Bike: true}},
		{[]string{"-transit"}, cmd.CommuteCmd{Transit: true}},
		{[]string{"-drive", "-walk"}, cmd.CommuteCmd{Drive: true, Walk: true}},
		{[]string{"-drive", "-walk", "-bike", "-transit"}, cmd.CommuteCmd{Drive: true, Walk: true, Bike: true, Transit: true}},
	}

	for idx, tt := range mTests {
		r, err := a.parseCommuteCmd(&conf, tt.args)
		if err != nil {
			t.Fatal(err)
		}

		if tt.expected.Drive != r.Drive {
			t.Fatalf("[%v] Unexpected 'Drive' parsed, expected=%v, got=%v", idx, tt.expected.Drive, r.Drive)
		} else if tt.expected.Walk != r.Walk {
			t.Fatalf("[%v] Unexpected 'Walk' parsed, expected=%v, got=%v", idx, tt.expected.Walk, r.Walk)
		} else if tt.expected.Bike != r.Bike {
			t.Fatalf("[%v] Unexpected 'Bike' parsed, expected=%v, got=%v", idx, tt.expected.Bike, r.Bike)
		} else if tt.expected.Transit != r.Transit {
			t.Fatalf("[%v] Unexpected 'Transit' parsed, expected=%v, got=%v", idx, tt.expected.Transit, r.Transit)
		}
	}
}

func TestArgParser_parseAddCmd(t *testing.T) {
	var a ArgParser
	var s MockStorageProvider

	tests := []struct {
		args     []string
		expected cmd.AddCmd
	}{
		{[]string{}, cmd.AddCmd{}},
		{[]string{"-name", "home"}, cmd.AddCmd{Name: "home"}},
		{[]string{"-location", "123 Main St."}, cmd.AddCmd{Value: "123 Main St."}},
		{[]string{"-name", "home", "-location", "123 Main St."}, cmd.AddCmd{Name: "home", Value: "123 Main St."}},
	}

	for idx, tt := range tests {
		r, err := a.parseAddCmd(&s, tt.args)
		if err != nil {
			t.Fatal(err)
		}

		if tt.expected.Name != r.Name {
			t.Fatalf("[%v] Unexpected 'Name' parsed, expected=%v, got=%v", idx, tt.expected.Name, r.Name)
		} else if tt.expected.Value != r.Value {
			t.Fatalf("[%v] Unexpected 'Value' parsed, expected=%v, got=%v", idx, tt.expected.Value, r.Value)
		} else if r.Store != &s {
			t.Fatalf("[%v] Unexpected Store, expected=%v, got=%v", idx, s, r.Store)
		}
	}
}

func testStringsEq(a, b []string) bool {
	if a == nil && b == nil {
		return true
	} else if a == nil || b == nil {
		return false
	} else if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
