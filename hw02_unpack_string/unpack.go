package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	runes := []rune(s)

	var result strings.Builder
	var next rune

	for i := 0; i < len(runes); i++ {
		current := runes[i]
		if i < len(runes)-1 {
			next = runes[i+1]
		} else {
			next = 0
		}

		// cant start with digit
		if i == 0 && unicode.IsDigit(current) {
			return "", ErrInvalidString
		}
		// two digits in a row
		if unicode.IsDigit(current) && unicode.IsDigit(next) {
			return "", ErrInvalidString
		}
		// if current is digit
		if unicode.IsDigit(next) {
			nextAsInt, _ := strconv.Atoi(string(next))
			result.WriteString(strings.Repeat(string(current), nextAsInt))
			continue
		}
		if !unicode.IsDigit(current) {
			result.WriteRune(current)
		}
	}

	return result.String(), nil
}
