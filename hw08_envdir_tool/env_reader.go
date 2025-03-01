package main

import (
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	// Place your code here
	dirContent, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	env := make(Environment, len(dirContent))

	for _, fileInfo := range dirContent {
		if fileInfo.IsDir() {
			continue
		}
		fileName := fileInfo.Name()

		if strings.Contains(fileName, "=") {
			continue
		}

		fileContent, err := os.ReadFile(dir + "/" + fileName)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %w", err)
		}
		firstLine, _, _ := strings.Cut(string(fileContent), "\n")
		firstLine = strings.TrimRight(firstLine, "\t ")
		firstLine = strings.ReplaceAll(firstLine, "\x00", "\n")

		info, err := fileInfo.Info()
		if err != nil {
			return nil, fmt.Errorf("error getting file info: %w", err)
		}
		env[fileName] = EnvValue{
			Value:      firstLine,
			NeedRemove: info.Size() == 0,
		}
	}

	return env, nil
}
