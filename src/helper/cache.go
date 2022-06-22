package helper

import (
	"fmt"
	"os"
	"strings"

	"github.com/ProtonMail/go-appdir"
	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/constants"
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
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = os.Remove(file.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
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
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		err = os.Remove(file.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}, nil
}
