package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(input string) (string, error) {
	if input == "" {
		return "", nil
	}
	var builder strings.Builder
	runesInput := []rune(input)

	pastIsSlash := false

	for idx := 0; idx < len(runesInput); idx++ {
		val := runesInput[idx]
		isLast := idx == len(runesInput)-1
		switch {
		case pastIsSlash:
			if unicode.IsDigit(val) || val == '\\' {
				writeString(&builder, runesInput, idx)
				pastIsSlash = false
			} else {
				return "", ErrInvalidString
			}
		case val == '\\':
			if isLast {
				return "", ErrInvalidString
			}
			pastIsSlash = true
		case unicode.IsDigit(val):
			if idx == 0 || (!isLast && unicode.IsDigit(runesInput[idx+1])) {
				return "", ErrInvalidString
			}
		default:
			writeString(&builder, runesInput, idx)
		}
	}
	return builder.String(), nil

}

func writeString(builder *strings.Builder, runeArr []rune, idx int) {
	if idx < len(runeArr)-1 && unicode.IsDigit(runeArr[idx+1]) {
		builder.WriteString(strings.Repeat(string(runeArr[idx]), int(runeArr[idx+1]-'0')))
	} else {
		builder.WriteString(string(runeArr[idx]))
	}
}

func main() {
	fmt.Println(Unpack("a3"))
}
