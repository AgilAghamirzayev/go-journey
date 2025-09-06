package main

import (
	"strings"
	"unicode"
)

func ModifyString(str string) string {

	trimmed := strings.TrimSpace(str)
	replaced := strings.Map(func(r rune) rune {

		if unicode.IsDigit(r) {
			return -1
		}
		return r
	}, trimmed)

	reversed := ""

	for _, r := range []rune(replaced) {
		reversed = string(r) + reversed
	}

	for _, r := range []rune(replaced) {
		reversed = string(r) + reversed
	}
	return reversed
}
