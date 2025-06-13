package utils

import (
	"testing"
)

func TestStringToInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"0", 0},
		{"-123", -123},
		{"abc", 0}, // invalid case
	}

	for _, test := range tests {
		result, err := StringToInt(test.input)
		if err != nil && test.expected != 0 {
			t.Errorf("StringToInt(%q) returned an error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("StringToInt(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}
func TestStringToInt64(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1234567890123456789", 1234567890123456789},
		{"0", 0},
		{"-1234567890", -1234567890},
		{"abc", 0}, // invalid case
	}

	for _, test := range tests {
		result, err := StringToInt64(test.input)
		if err != nil && test.expected != 0 {
			t.Errorf("StringToInt64(%q) returned an error: %v", test.input, err)
		}
		if result != test.expected {
			t.Errorf("StringToInt64(%q) = %d; want %d", test.input, result, test.expected)
		}
	}
}
