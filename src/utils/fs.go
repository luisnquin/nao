package utils

import "os"

func IsDirectory(path string) (bool, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fi.IsDir(), nil
}
