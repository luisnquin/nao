package utils

import (
	"os"
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
