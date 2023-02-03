package utils

import (
	"os"

	"github.com/agnivade/levenshtein"
)

func Contains[T comparable](slice []T, target T) bool {
	for _, el := range slice {
		if el == target {
			return true
		}
	}

	return false
}

func IsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}

func Ptr[T any](v T) *T {
	return &v
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
