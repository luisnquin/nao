package data

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/luisnquin/nao/src/config"
	"github.com/spf13/cobra"
)

func New() *Box {
	var box Box

	err := os.MkdirAll(config.App.Paths.DataDir, os.ModePerm)
	cobra.CheckErr(err)

	var f *os.File

	if _, err = os.Stat(config.App.Paths.DataFile); errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(config.App.Paths.DataFile)
	} else {
		f, err = os.Open(config.App.Paths.DataFile)
	}

	cobra.CheckErr(err)

	defer f.Close()

	err = json.NewDecoder(f).Decode(&box.box)
	if err != nil && !errors.Is(err, io.EOF) {
		panic(err)
	}

	if box.box.NaoSet == nil {
		box.box.NaoSet = make(map[string]Note, 0)
	}

	return &box
}

func JustLoadBox() BoxData {
	f, err := os.Open(config.App.Paths.DataFile)
	if err != nil {
		panic(err)
	}

	var data BoxData

	err = json.NewDecoder(f).Decode(&data)
	if err != nil {
		panic(err)
	}

	return data
}
