package helpers

import "testing"

func TestTrimSpaceToupper(t *testing.T) {
	type test struct {
		name     string
		input    string
		expected string
	}

	tests := []test{
		{
			"basic test: all lowercase with leading ending space",
			"  this is a test   ",
			"THIS IS A TEST",
		},
		{
			"basic test: all uppercase with leading ending space",
			"  THIS IS A TEST   ",
			"THIS IS A TEST",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := TrimSpaceToUpper(test.input)
			if got != test.expected {
				t.Errorf("got %s\n want %s\n", got, test.expected)
			}
		})
	}
}
