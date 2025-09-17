package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	expected := make(Environment)
	expected["BAR"] = EnvValue{Value: "bar"}
	expected["EMPTY"] = EnvValue{Value: ""}
	expected["FOO"] = EnvValue{Value: "   foo\nwith new line"}
	expected["HELLO"] = EnvValue{Value: "\"hello\""}
	expected["UNSET"] = EnvValue{Value: "", NeedRemove: true}

	t.Run("Read dir standart", func(t *testing.T) {
		result, err := ReadDir("./testdata/env")
		if err != nil {
			t.Fatal(err)
		}
		require.Equal(t, expected, result)
	})
}
