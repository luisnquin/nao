package utils

import (
	"encoding/json"
	"fmt"
)

const (
	MB = 1 << 20
	GB = 1 << 30
	KB = 1 << 10
)

// Returns an human-readable representation of the size of the passed JSON value.
func GetHumanReadableSize(v any) string {
	return SizeToStorageUnits(int64(GetSize(v)))
}

// Returns the size of the passed JSON value.
func GetSize(v any) int {
	content, _ := json.Marshal(v)

	return len(content)
}

func SizeToStorageUnits[T int64 | int | float64](n T) string {
	switch {
	case n > GB:
		return fmt.Sprintf("%.2f%s", float64(n)/(GB), "GB")

	case n > MB:
		return fmt.Sprintf("%.2f%s", float64(n)/(MB), "MB")

	default:
		return fmt.Sprintf("%.2f%s", float64(n)/(KB), "KB")
	}
}
