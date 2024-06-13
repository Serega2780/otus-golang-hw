package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	init := []rune(str)
	var sb strings.Builder
	backSlash := false
	for i, v := range init {
		_, err := check(i, v, init)
		if err != nil {
			return "", err
		}

		switch {
		case v == 92:
			if backSlash {
				write(&sb, v)
				backSlash = false
			} else {
				backSlash = true
			}
		case unicode.IsDigit(v):
			count, err := strconv.Atoi(string(v))
			if err != nil {
				return "", ErrInvalidString
			}
			if backSlash {
				write(&sb, v)
				backSlash = false
			} else {
				err = repeat(&sb, init[i-1], count-1)
				if err != nil {
					return "", err
				}
			}
		default:
			if backSlash {
				return "", ErrInvalidString
			}
			write(&sb, v)
		}
	}
	if backSlash {
		return "", ErrInvalidString
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

func repeat(sb *strings.Builder, symbol rune, count int) error {
	if count < 0 {
		err := trimLastRune(sb)
		if err != nil {
			return err
		}
	} else {
		sb.WriteString(strings.Repeat(string(symbol), count))
	}

	return nil
}

func trimLastRune(sb *strings.Builder) error {
	err, size := utf8.DecodeLastRuneInString(sb.String())
	if err == utf8.RuneError && (size == 0 || size == 1) {
		return ErrInvalidString
	}
	str := sb.String()
	sb.Reset()
	sb.WriteString(str[:len(str)-size])
	return nil
}
