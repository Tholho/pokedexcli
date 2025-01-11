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
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "    ",
			expected: []string{},
		},
		{
			input:    "  h e l l o  world  ",
			expected: []string{"h", "e", "l", "l", "o", "world"},
		},
		{
			input:    "  %%  ..  ",
			expected: []string{"%%", ".."},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("length don't match")
			t.Fail()
		}
		// Check the length of the actual slice
		// if they don't match, use t.Errorf to print an error message
		// and fail the test
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			t.Log("expecting: ", expectedWord, "\ngot: ", word)
			if word != expectedWord {
				t.Errorf("words dont match")
				t.Fail()
			}
			// Check each word in the slice
			// if they don't match, use t.Errorf to print an error message
			// and fail the test
		}
	}

}
