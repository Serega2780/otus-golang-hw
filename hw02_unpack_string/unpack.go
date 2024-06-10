package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	init := []rune(str)
	var sb strings.Builder
	bcsl := false
	for i, v := range init {
		_, err := check(i, v, init)
		if err != nil {
			return "", err
		}

		switch {
		case v == 92:
			if bcsl {
				write(&sb, v)
				bcsl = false
			} else {
				bcsl = true
			}
		case unicode.IsDigit(v):
			count, err := strconv.Atoi(string(v))
			if err != nil {
				return "", ErrInvalidString
			}
			if bcsl {
				write(&sb, v)
				bcsl = false
			} else {
				repeat(&sb, init[i-1], count-1)
			}
		default:
			write(&sb, v)
		}
	}

	return sb.String(), nil
}

func check(i int, symbol rune, init []rune) (bool, error) {
	if i == 0 && unicode.IsDigit(symbol) {
		return false, ErrInvalidString
	}
	if unicode.IsDigit(symbol) && unicode.IsDigit(init[i-1]) && init[i-2] != 92 {
		return false, ErrInvalidString
	}

	return true, nil
}

func write(sb *strings.Builder, symbol rune) {
	sb.WriteRune(symbol)
}

func repeat(sb *strings.Builder, symbol rune, count int) {
	if count < 0 {
		str := sb.String()
		str = str[:len(str)-1]
		sb.Reset()
		sb.WriteString(str)
	} else {
		sb.WriteString(strings.Repeat(string(symbol), count))
	}
}
