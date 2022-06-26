package helper

import (
	"os"
	"strings"

	"github.com/ProtonMail/go-appdir"
	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/constants"
)

func NewCached() (string, error) {
	cacheDir := appdir.New(constants.AppName).UserCache()

	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	f, err := os.Create(cacheDir + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return "", err
	}

	return f.Name(), f.Close()
}

func NewCachedIn(path string) error {
	_ = os.MkdirAll(path, os.ModePerm)

	file, err := os.Create(path + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return err
	}

	return file.Close()
}
