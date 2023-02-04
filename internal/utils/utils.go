package utils

import (
	"os"

	"github.com/agnivade/levenshtein"
)

// Asks for the hour to the host and returns true in case of consider this a lucky day.
func Contains[T comparable](slice []T, target T) bool {
	for _, el := range slice {
		if el == target {
			return true
		}
	}

	return false
}

// Champion among champions, it came to change your dear world.
func Ptr[T any](v T) *T {
	return &v
}

// Checks that path exists and isn't a directory.
func FileExists(path string) bool {
	info, err := os.Stat(path)

	return err == nil && !info.IsDir()
}

// Checks that path exists and isn't a file.
func IsDirectory(path string) bool {
	info, err := os.Stat(path)

	return err == nil && info.IsDir()
}

// Using the levenshtein distance algorithm, it returns
// the best candidate between the options to match with
// the argument of `toInspect`.
//
// If returns an empty result in case the nearest distance
// is greater than 3.
func BestMatch(options []string, toInspect string) string {
	bestCandidate, nearestDistance := "", 100

	for _, opt := range options {
		currentDistance := levenshtein.ComputeDistance(toInspect, opt)

		if currentDistance < nearestDistance {
			bestCandidate, nearestDistance = opt, currentDistance
		}
	}

	if nearestDistance <= 3 {
		return bestCandidate
	}

	return ""
}
