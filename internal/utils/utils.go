package utils

import (
	"os"
)

func Contains[T int | string](slice []T, target T) bool {
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
