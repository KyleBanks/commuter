package cli

import (
	"bytes"
	"os"
	"testing"
)

func TestNewStdout(t *testing.T) {
	s := NewStdout()

	if s.Writer != os.Stdout {
		t.Fatalf("Unexpected Writer, expected=os.Stdout, got=%v", s.Writer)
	}
}

func TestStdout_Indicate(t *testing.T) {
	tests := []struct {
		msg      string
		args     []interface{}
		expected string
	}{
		{"Hey", nil, "Hey\n"},
		{"Hey %v", []interface{}{"There"}, "Hey There\n"},
		{"I am %v and I am %v", []interface{}{"A Test", true}, "I am A Test and I am true\n"},
	}

	for idx, tt := range tests {
		var b bytes.Buffer
		s := Stdout{
			Writer: &b,
		}

		s.Indicate(tt.msg, tt.args...)
		out := b.String()

		if out != tt.expected {
			t.Fatalf("[#%v] Unexpected output, expected=%v, got=%v", idx, tt.expected, out)
		}
	}
}

func TestNewStdin(t *testing.T) {
	s := NewStdin()
	if s.Scanner == nil {
		t.Fatalf("Unexpected nil Scanner")
	}
}
