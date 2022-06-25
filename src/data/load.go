package data

import (
	"encoding/json"
	"errors"
	"fmt"
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

	err = json.NewDecoder(f).Decode(&box.data)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if box.data.NaoSet == nil {
		box.data.NaoSet = make(map[string]Set, 0)
	}

	return &box
}
