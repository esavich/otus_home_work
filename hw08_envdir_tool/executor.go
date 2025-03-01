package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	// Place your code here.

	command := exec.Command(cmd[0], cmd[1:]...) //nolint

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	currEnv := os.Environ()

	// append current environment variables, without marked to remove
	for _, value := range currEnv {
		envKey, _, _ := strings.Cut(value, "=")
		if _, ok := env[envKey]; !ok || !env[envKey].NeedRemove {
			command.Env = append(command.Env, value)
		}
	}

	// append new environment variables
	for key, value := range env {
		if !value.NeedRemove {
			command.Env = append(command.Env, fmt.Sprintf("%s=%s", key, value.Value))
		}
	}

	err := command.Run()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return command.ProcessState.ExitCode()
}
