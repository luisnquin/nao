package utils

import (
	"encoding/json"
	"fmt"
)

// Returns an human-readable representation of the size of the passed JSON value.
func GetHumanReadableSize(v any) string {
	return SizeToStorageUnits(int64(GetSize(v)))
}

// Returns the size of the passed JSON value
func GetSize(v any) int {
	content, _ := json.Marshal(v)

	return len(content)
}

func SizeToStorageUnits[T int64 | int | float64](n T) string {
	switch {
	case n > 1000000:
		return fmt.Sprintf("%.2f%s", float64(n)/(1<<20), "MB")

	case n > 1_000_000_000:
		return fmt.Sprintf("%.2f%s", float64(n)/(1<<30), "GB")

	default:
		return fmt.Sprintf("%.2f%s", float64(n)/(1<<10), "KB")
	}
}
