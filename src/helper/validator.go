package helper

import (
	"regexp"
)

func EnsureTagIsValid(tag string) bool {
	return regexp.MustCompile(`^[A-z0-9\@\_\-]+$`).MatchString(tag)
}
