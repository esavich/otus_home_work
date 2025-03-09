package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	// Place your code here

	originalStdout := os.Stdout
	r, w, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = w
	// also test that existing environment variables are handled
	err = os.Setenv("MUST_BE_REMOVED", "MBR")
	require.NoError(t, err)
	err = os.Setenv("MUST_EXIST_FROM_OS", "ME")
	require.NoError(t, err)

	env := Environment{
		"FOO":             {Value: "foo", NeedRemove: false},
		"BAR":             {Value: "bar", NeedRemove: false},
		"MUST_BE_REMOVED": {Value: "", NeedRemove: true},
	}

	cmd := []string{"/bin/sh", "-c", "echo $FOO $BAR $MUST_EXIST_FROM_OS $MUST_BE_REMOVED"}

	exitCode := RunCmd(cmd, env)
	require.Equal(t, 0, exitCode)

	err = w.Close()
	require.NoError(t, err)

	os.Stdout = originalStdout

	var buf bytes.Buffer
	_, err = r.WriteTo(&buf)
	require.NoError(t, err)

	output := buf.String()

	require.Equal(t, "foo bar ME\n", output)
}
