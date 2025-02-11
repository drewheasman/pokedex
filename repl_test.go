package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "",
			expected: []string{},
		},
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "HELLO, world!",
			expected: []string{"hello,", "world!"},
		},
		{
			input:    "Charmander Bulbasaur PIKACHU",
			expected: []string{"charmander", "bulbasaur", "pikachu"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(c.expected) != len(actual) {
			t.Errorf("expected word count %d (%v), actual %d (%v)", len(c.expected), c.expected, len(actual), actual)
		}
		for i := range actual {
			expectedWord := c.expected[i]
			word := actual[i]

			if word != expectedWord {
				t.Errorf("expected word %s, actual %s", expectedWord, word)
			}
		}
	}
}
