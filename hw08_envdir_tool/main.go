package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	args := os.Args
	fmt.Println(args)

	if len(args) < 3 {
		fmt.Println("to few arguments")
		fmt.Println("Usage: go-envdir path/to/env/dir command arg1 arg2 ...")
		os.Exit(1)
	}

	envDir := args[1]
	cmd := args[2:]

	fmt.Printf("envdir: %s, cmd and args %s\n", envDir, cmd)

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(env)
}
