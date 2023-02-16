package utils

import (
	"strings"
	"unicode"
)

func ToPascalCase(s string) string {
	return toTitleCase(s, true)
}

func ToCamelCase(s string) string {
	return toTitleCase(s, false)
}

func toTitleCase(s string, capitalizeFirst bool) string {
	var result strings.Builder
	result.Grow(len(s))

	for _, r := range s {
		if r == ' ' {
			capitalizeFirst = true
		} else if capitalizeFirst {
			result.WriteRune(unicode.ToUpper(r))
			capitalizeFirst = false
		} else {
			result.WriteRune(unicode.ToLower(r))
		}
	}

	return result.String()
}
