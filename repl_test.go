package main

import "testing"

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
			input:    "quick brown fox",
			expected: []string{"quick", "brown", "fox"},
		},
		{
			input:    "    leading spaces",
			expected: []string{"leading", "spaces"},
		},
		{
			input:    "trailing spaces    ",
			expected: []string{"trailing", "spaces"},
		},
		{
			input:    "  multiple    spaces  between  words  ",
			expected: []string{"multiple", "spaces", "between", "words"},
		},
		{
			input:    "",
			expected: []string{},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		actualLength := len(actual)
		expectedLength := len(c.expected)
		if actualLength != expectedLength {
			t.Errorf("Length of slices are not the same!\nGot: %d\nExpected: %d", actualLength, expectedLength)
			t.Fail()
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("Cannot find word!\nGot: %s\nExpected: %s", word, expectedWord)
				t.Fail()
			}
		}
	}
}
