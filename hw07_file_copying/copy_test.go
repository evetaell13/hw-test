package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopySuccess(t *testing.T) {
	tests := []struct {
		from, to      string
		limit, offset int64
		expected      error
	}{
		{from: "testdata/input.txt", to: "out.txt", limit: 0, offset: 0, expected: nil},
		{from: "testdata/input.txt", to: "out.txt", limit: 10, offset: 0, expected: nil},
		{from: "testdata/input.txt", to: "out.txt", limit: 1000, offset: 0, expected: nil},
		{from: "testdata/input.txt", to: "out.txt", limit: 10000, offset: 0, expected: nil},
		{from: "testdata/input.txt", to: "out.txt", limit: 1000, offset: 10, expected: nil},
	}

	for _, tc := range tests {
		tc := tc
		t.Run("copy without errors", func(t *testing.T) {
			err := Copy(tc.from, tc.to, tc.offset, tc.limit)
			require.NoError(t, err)
		})
		os.Remove(tc.to)
	}
}

func TestCopyErrorOffset(t *testing.T) {
	tc := struct {
		from, to      string
		limit, offset int64
		expected      error
	}{
		from:     "testdata/input.txt",
		to:       "out.txt",
		limit:    10,
		offset:   10000,
		expected: ErrOffsetExceedsFileSize,
	}

	t.Run("copy error offset", func(t *testing.T) {
		err := Copy(tc.from, tc.to, tc.offset, tc.limit)
		require.Equal(t, tc.expected, err)
	})
}

func TestCopyErrorFileSize(t *testing.T) {
	tc := struct {
		from, to      string
		limit, offset int64
		expected      error
	}{
		from:     "/dev/urandom",
		to:       "out.txt",
		limit:    1000,
		offset:   0,
		expected: ErrUnsupportedFile,
	}

	t.Run("copy error size", func(t *testing.T) {
		err := Copy(tc.from, tc.to, tc.offset, tc.limit)
		require.Equal(t, tc.expected, err)
	})
}

func TestCopyErrorOpenFile(t *testing.T) {
	tc := struct {
		from, to      string
		limit, offset int64
		expected      error
	}{
		from:     "testdata/input123.txt",
		to:       "out.txt",
		limit:    1000,
		offset:   0,
		expected: ErrMessageOpenFile,
	}

	t.Run("copy error size", func(t *testing.T) {
		err := Copy(tc.from, tc.to, tc.offset, tc.limit)
		require.Equal(t, tc.expected, err)
	})
}
