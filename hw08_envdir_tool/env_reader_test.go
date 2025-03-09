package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	// Place your code here
	files := map[string]string{
		"FIRST":             "1st",
		"SECOND":            "2nd",
		"INVALID=NAME":      "no matter",
		"=INVALID_NAME":     "12345",
		"MULTILINE":         "firstline\nsecondline",
		"TOBETRIMMED":       "something \t",
		"PRESERVELEFTSPACE": "  something",
		"NULLSASNEWLINE":    "textwith\x00newline",
		"UNSET":             "",
		"EMPTY":             "\n",
	}

	t.Run("test current testdata", func(t *testing.T) {
		env, err := ReadDir("./testdata/env")
		require.NoError(t, err)

		expected := Environment{
			"BAR":   {Value: "bar", NeedRemove: false},
			"EMPTY": {Value: "", NeedRemove: false},
			"FOO":   {Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": {Value: "\"hello\"", NeedRemove: false},
			"UNSET": {Value: "", NeedRemove: true},
		}

		require.Equal(t, expected, env)
	})

	t.Run("empty directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		env, err := ReadDir(tmpDir)
		require.NoError(t, err)

		require.Empty(t, env)
	})

	t.Run("some files in  directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		for name, content := range files {
			err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0o666)
			require.NoError(t, err)
		}

		env, err := ReadDir(tmpDir)
		require.NoError(t, err)

		expected := Environment{
			"FIRST":             {Value: "1st", NeedRemove: false},
			"SECOND":            {Value: "2nd", NeedRemove: false},
			"MULTILINE":         {Value: "firstline", NeedRemove: false},
			"TOBETRIMMED":       {Value: "something", NeedRemove: false},
			"PRESERVELEFTSPACE": {Value: "  something", NeedRemove: false},
			"NULLSASNEWLINE":    {Value: "textwith\nnewline", NeedRemove: false},
			"UNSET":             {Value: "", NeedRemove: true},
			"EMPTY":             {Value: "", NeedRemove: false},
		}

		require.Equal(t, expected, env)
	})

	t.Run("some files in sub directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		subdirPath := filepath.Join(tmpDir, "subdir")

		err := os.Mkdir(subdirPath, 0o777)
		require.NoError(t, err)

		// place files in subdir
		for name, content := range files {
			err := os.WriteFile(filepath.Join(subdirPath, name), []byte(content), 0o666)
			require.NoError(t, err)
		}
		// read tmpDir
		env, err := ReadDir(tmpDir)
		require.NoError(t, err)

		// subdir should be ignored
		require.Empty(t, env)
	})
}
