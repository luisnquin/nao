package helper

import (
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/config"
)

func NewCached() (string, error) {
	err := os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(config.App.Paths.CacheDir + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return "", err
	}

	return f.Name(), f.Close()
}

func NewCachedIn(path string) error {
	_ = os.MkdirAll(config.App.Paths.CacheDir, os.ModePerm)

	file, err := os.Create(path + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return err
	}

	return file.Close()
}
