package utils

import "strings"

func Contains(args []string, target string) bool {
	for _, arg := range args {
		if arg == target {
			return true
		}
	}

	return false
}

func PrettyJoin(v []string) string {
	return strings.Join(v[:len(v)-1], ", ") + " and " + v[len(v)-1]
}
