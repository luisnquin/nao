package utils

import (
	"strings"
	"unicode"
)

func ToCamelCase(s string) string {
	var result strings.Builder
	result.Grow(len(s))

	capitalize := false

	for _, r := range s {
		if r == ' ' {
			capitalize = true
		} else if capitalize {
			result.WriteRune(unicode.ToUpper(r))
			capitalize = false
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}

	return result.String()
}
