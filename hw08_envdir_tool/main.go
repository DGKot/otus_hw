package main

import (
	"os"
)

func main() {
	// Place your code here.
	path := os.Args[1]
	env, err := ReadDir(path)
	if err != nil {
		return
	}
	code := RunCmd(os.Args[2:], env)
	os.Exit(code)
}
