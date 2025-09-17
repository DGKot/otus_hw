package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestCopy(t *testing.T) {
	testCases := []struct {
		name     string
		text     string
		expected string
		limit    int64
	}{
		{
			name:     "Полное копирование",
			text:     "Test text",
			expected: "Test text",
			limit:    int64(len("Test text")),
		},
		{
			name:     "Частичное копирование",
			text:     "Test text",
			expected: "Test",
			limit:    int64(len("Test")),
		},
		{
			name:     "Пустая строка",
			text:     "",
			expected: "",
			limit:    0,
		},
	}
	for _, testCase := range testCases {
		from := strings.NewReader(testCase.text)
		to := &bytes.Buffer{}
		t.Run(testCase.name, func(t *testing.T) {
			err := CopyCore(to, from, testCase.limit)
			if err != nil {
				t.Fatal(err)
			}
			res := to.String()
			if res != testCase.expected {
				t.Errorf("got %s, expected %s", res, testCase.expected)
			}
		})
	}

	// Place your code here.
}
