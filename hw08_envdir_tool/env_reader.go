package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool // if file is empty = true else = false
}

var (
	ErrEmptyPath = errors.New("empty path of dir")
	ErrEmptyDir  = errors.New("empty dir")
)

func (e *Environment) ToSliceString() []string {
	res := []string{}
	for k, v := range *e {
		res = append(res, fmt.Sprintf("%s=%s", k, v.Value))
	}
	return res
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	if dir == "" {
		return nil, ErrEmptyPath
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, ErrEmptyDir
	}
	envs := make(Environment)
	for _, obj := range entries {
		if obj.IsDir() || strings.Contains(obj.Name(), "=") {
			continue
		}
		file, err := os.Open(dir + "/" + obj.Name())
		if err != nil {
			return nil, err
		}
		fileInfo, err := file.Stat()
		if err != nil {
			return nil, err
		}
		os.Unsetenv(obj.Name())
		if fileInfo.Size() == 0 {
			envs[obj.Name()] = EnvValue{NeedRemove: true}
			file.Close()
			continue
		}
		scanner := bufio.NewScanner(file)
		if scanner.Scan() {
			envs[obj.Name()] = EnvValue{Value: strings.ReplaceAll(strings.TrimRight(scanner.Text(), "\t "), "\x00", "\n")}
		}

		file.Close()
	}

	return envs, nil
}
