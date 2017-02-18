package main

import (
	"testing"
	"time"
)

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
