package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "🙃0", expected: ""},
		{input: "🙃", expected: "🙃"},
		{input: "aaф0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "abc\n", expected: "abc\n"},
		{input: "日本語", expected: "日本語"},
		{input: "日3本\n語", expected: "日日日本\n語"},

		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: `qwe\\\3\\`, expected: `qwe\3\`},
		{input: `s`, expected: `s`},
		{input: `w1`, expected: `w`},
		{input: `a\\0b`, expected: `ab`},
		{input: `\\3\\4\\5`, expected: `\\\\\\\\\\\\`},
		{input: `日3本\2語`, expected: "日日日本2語"},
		{input: `🙃4本\\🙃2語`, expected: `🙃🙃🙃🙃本\🙃🙃語`},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)

			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", `qwe\\\3\`, `\ `, `日3本\n語`}
	for _, tc := range invalidStrings {
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
