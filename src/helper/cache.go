package helper

import (
	"os"
	"strings"

	"github.com/ProtonMail/go-appdir"
	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/constants"
	"github.com/spf13/cobra"
)

func NewCached() (*os.File, func(), error) {
	cacheDir := appdir.New(constants.AppName).UserCache()

	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Create(cacheDir + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		err = file.Close()
		cobra.CheckErr(err)

		err = os.Remove(file.Name())
		cobra.CheckErr(err)
	}, nil
}

func NewCachedIn(path string) (*os.File, func(), error) {
	_ = os.MkdirAll(path, os.ModePerm)

	file, err := os.Create(path + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		return nil, nil, err
	}

	return file, func() {
		err = file.Close()
		cobra.CheckErr(err)

		err = os.Remove(file.Name())
		cobra.CheckErr(err)
	}, nil
}
