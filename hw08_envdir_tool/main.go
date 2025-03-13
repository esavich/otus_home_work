package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	args := os.Args

	if len(args) < 3 {
		fmt.Println("to few arguments")
		fmt.Println("Usage: go-envdir path/to/env/dir command arg1 arg2 ...")
		os.Exit(1)
	}

	envDir := args[1]
	cmd := args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		os.Exit(1)
	}

	returnCode := RunCmd(cmd, env)

	os.Exit(returnCode)
}
