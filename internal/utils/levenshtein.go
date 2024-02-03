package utils

import "github.com/agnivade/levenshtein"

const maximumLevenshteinDistance = 3

// Using the levenshtein distance algorithm. It returns
// the best candidate between the options to match with
// the target.
//
// If returns an empty result in case the nearest distance
// is greater than 3.
func CalculateNearestString(candidates []string, target string) string {
	bestCandidate, shorterDistance := "", 100

	for _, opt := range candidates {
		currentDistance := levenshtein.ComputeDistance(target, opt)

		if currentDistance < shorterDistance {
			bestCandidate, shorterDistance = opt, currentDistance
		}
	}

	if shorterDistance > maximumLevenshteinDistance {
		return ""
	}

	return bestCandidate
}
