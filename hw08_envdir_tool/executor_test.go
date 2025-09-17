package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Place your code here
	env := make(Environment)
	env["foo"] = EnvValue{Value: "foo value"}
	env["bar"] = EnvValue{Value: "bar value"}
	cmd := []string{"sh", "-c", "echo \"foo is ${foo}\nbar is ${bar} arg1 is $0 arg2 is $1\"", "5", "10"}
	expectedString := "foo is foo value\nbar is bar value arg1 is 5 arg2 is 10\n"
	resultString := customOutput(func() {
		code := RunCmd(cmd, env)
		require.Equal(t, 0, code)
	})
	require.Equal(t, expectedString, resultString)
}

func customOutput(f func()) string {
	sysOut := os.Stdout
	sysErr := os.Stderr

	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	f()

	_ = w.Close()
	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	os.Stdout = sysOut
	os.Stderr = sysErr

	return buf.String()
}
