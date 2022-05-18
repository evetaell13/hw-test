package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "mm4n2e0", expected: "mmmmmnn"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "aaa0b0", expected: "aa"},
		{input: "id2qd", expected: "iddqd"},
		{input: "idkk0f1aa0", expected: "idkfa"},
		{input: "m0o0r0", expected: ""},
		// uncomment if task with asterisk completed
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackStringKirilic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "г4м3с2р1", expected: "ггггмммсср"},
		{input: "к2м", expected: "ккм"},
		{input: "оти1с", expected: "отис"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackWithBigLetterString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "AAB5C2N", expected: "AABBBBBCCN"},
		{input: "A0NN2B0", expected: "NNN"},
		{input: "C\n2c3", expected: "C\n\nccc"},
		{input: "L0k0M0", expected: ""},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "1", `rw\ne`, "aaa00", "A44B", "лл3л67"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}
