package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	// Place your code here.

	t.Run("negative offset", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/output.txt", -1, 0)

		require.ErrorIs(t, err, ErrNegativeOffset)
	})

	t.Run("negative limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "testdata/output.txt", 0, -1)

		require.ErrorIs(t, err, ErrorNegativeLimit)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		input, err := os.CreateTemp(os.TempDir(), "input")
		if err != nil {
			t.Error(err)
		}
		require.NoError(t, err)
		defer removeTmp(t, input)

		err = Copy(input.Name(), "testdata/output.txt", 1000, 0)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("testdata", func(t *testing.T) {
		tests := []struct {
			offset int64
			limit  int64
		}{
			{0, 0},
			{0, 10},
			{0, 1000},
			{0, 10000},
			{100, 1000},
			{6000, 1000},
		}
		for _, tc := range tests {
			tc := tc
			fileName := buildName(tc.offset, tc.limit)

			t.Run(fileName, func(t *testing.T) {
				output, err := os.CreateTemp(os.TempDir(), "output")
				if err != nil {
					t.Error(err)
				}
				require.NoError(t, err)
				defer removeTmp(t, output)

				err = Copy("testdata/input.txt", output.Name(), tc.offset, tc.limit)
				require.NoError(t, err)

				result, err := os.ReadFile(output.Name())
				require.NoError(t, err)

				expected, err := os.ReadFile(fileName)
				require.NoError(t, err)

				require.Equal(t, expected, result)
			})
		}
	})
}

func buildName(o int64, l int64) string {
	return fmt.Sprintf("testdata/out_offset%d_limit%d.txt", o, l)
}

func removeTmp(t *testing.T, input *os.File) {
	t.Helper()
	err := os.Remove(input.Name())
	require.NoError(t, err)
	if err != nil {
		t.Error(err)
	}
}
