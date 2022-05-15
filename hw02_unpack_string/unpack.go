package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var backslash bool
	sr := []rune(str)
	var count int
	var strResult string

	for i, item := range sr {
		if i == 0 && unicode.IsDigit(item) {
			return "", ErrInvalidString
		}

		if unicode.IsDigit(item) && unicode.IsDigit(sr[i-1]) && sr[i-2] != '\\' {
			return "", ErrInvalidString
		}

		if item == '\\' && !backslash {
			backslash = true
			continue
		}

		if backslash && unicode.IsLetter(item) {
			return "", ErrInvalidString
		}

		if backslash {
			strResult += string(item)
			backslash = false
			continue
		}

		if unicode.IsDigit(item) {
			count = int(item - '0')
			if count == 0 {
				strResult = strResult[:len(strResult)-1]
				continue
			}
			strResult += strings.Repeat(string(sr[i-1]), count-1)
			continue
		}
		strResult += string(item)
	}

	return strResult, nil
}
