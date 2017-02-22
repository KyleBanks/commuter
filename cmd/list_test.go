package cmd

import (
	"strings"
	"testing"
)

func TestListCmd_Run(t *testing.T) {
	tests := []struct {
		locs   map[string]string
		expect [][2]string
	}{
		{
			map[string]string{"aname": "value1", "bname": "value2"},
			[][2]string{
				{"aname", "value1"},
				{"bname", "value2"},
			},
		},
		{
			map[string]string{"bname": "value1", "aname": "value2"},
			[][2]string{
				{"aname", "value2"},
				{"bname", "value1"},
			},
		},
		{
			map[string]string{"bname": "value1", "aname": "value2", "default": "defaultval"},
			[][2]string{
				{"default", "defaultval"},
				{"aname", "value2"},
				{"bname", "value1"},
			},
		},
	}

	for idx, tt := range tests {
		var l ListCmd
		var m mockIndicator
		conf := Configuration{Locations: tt.locs}

		if err := l.Run(&conf, &m); err != nil {
			t.Fatal(err)
		}

		if len(m.out) != len(tt.expect) {
			t.Fatalf("[#%v] Unexpected number of output lines, expected=%v, got=%v", idx, len(tt.expect), len(m.out))
		}

		for line, ex := range tt.expect {
			out := strings.TrimSpace(m.out[line])
			expect := ex[0] + ": " + ex[1]

			if out != expect {
				t.Fatalf("[#%v] [Line %v] Unexpected output, expected=%v, got=%v", idx, line, expect, out)
			}
		}
	}
}

func TestListCmd_Validate(t *testing.T) {
	tests := []struct {
		conf Configuration
	}{
		{Configuration{}},
		{Configuration{Locations: make(map[string]string)}},
		{Configuration{Locations: map[string]string{"name": "value"}}},
	}

	for _, tt := range tests {
		var l ListCmd
		if err := l.Validate(&tt.conf); err != nil {
			t.Fatal(err)
		}
	}
}
