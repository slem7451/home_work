package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	strRune := []rune(str)
	var res strings.Builder

	if len(strRune) > 0 && unicode.IsDigit(strRune[0]) {
		return "", ErrInvalidString
	}

	for i := 0; i < len(strRune); i++ {
		switch {
		case string(strRune[i]) == `\`:
			if i+1 < len(strRune) && (unicode.IsDigit(strRune[i+1]) || string(strRune[i+1]) == `\`) { //nolint:nestif
				if i+2 < len(strRune) && unicode.IsDigit(strRune[i+2]) {
					if i+3 < len(strRune) && unicode.IsDigit(strRune[i+3]) {
						return "", ErrInvalidString
					}

					digit, _ := strconv.Atoi(string(strRune[i+2]))
					res.WriteString(strings.Repeat(string(strRune[i+1]), digit))
				} else {
					res.WriteString(string(strRune[i+1]))
				}
				i += 2
			} else {
				return "", ErrInvalidString
			}
		case i+1 < len(strRune) && unicode.IsDigit(strRune[i+1]):
			if i+2 < len(strRune) && unicode.IsDigit(strRune[i+2]) {
				return "", ErrInvalidString
			}

			digit, _ := strconv.Atoi(string(strRune[i+1]))
			res.WriteString(strings.Repeat(string(strRune[i]), digit))
			i++
		default:
			res.WriteString(string(strRune[i]))
		}
	}

	return res.String(), nil
}
