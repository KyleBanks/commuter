package cmd

import (
	"errors"
	"testing"
)

func TestAddCmd_Run(t *testing.T) {
	tests := []struct {
		conf *Configuration

		name  string
		value string
	}{
		{&Configuration{Locations: make(map[string]string)}, "name", "value"},
		{&Configuration{Locations: map[string]string{"name": "value1"}}, "name", "value2"},
	}

	for idx, tt := range tests {
		m := mockStorageProvider{
			saveFn: func(i interface{}) error {
				if i != tt.conf {
					t.Fatalf("[#%v] Unexpected parameter to Save, got=%v", idx, i)
				}
				return nil
			},
		}
		a := AddCmd{Name: tt.name, Value: tt.value, Store: &m}

		if err := a.Run(tt.conf, nil); err != nil {
			t.Fatal(err)
		}

		if tt.conf.Locations[tt.name] != tt.value {
			t.Fatalf("[#%v] Unexpected value stored, expected=%v, got=%v", idx, tt.value, tt.conf.Locations[tt.name])
		}
	}

	// Negative
	{
		testErr := errors.New("mock err")
		var m mockStorageProvider
		m.saveFn = func(i interface{}) error {
			return testErr
		}

		conf := Configuration{Locations: make(map[string]string)}
		a := AddCmd{Name: "name", Value: "value", Store: &m}
		if err := a.Run(&conf, nil); err != testErr {
			t.Fatalf("Unexpected error, expected=%v, got=%v", testErr, err)
		}
	}
}

func TestAddCmd_Validate(t *testing.T) {
	tests := []struct {
		name  string
		value string
		err   error
	}{
		{"name", "value", nil},
		{"", "value", ErrAddNameMissing},
		{"name", "", ErrAddLocationMissing},
	}

	for idx, tt := range tests {
		a := AddCmd{Name: tt.name, Value: tt.value}

		if err := a.Validate(nil); err != tt.err {
			t.Fatalf("[#%v] Unexpected error, expected=%v, got=%v", idx, tt.err, err)
		}
	}
}
