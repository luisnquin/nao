package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/luisnquin/nao/src/config"
)

func NewUserBox() *Box {
	var box Box

	dataDir := config.App.Dirs.UserData()

	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var f *os.File

	if _, err = os.Stat(dataDir + "/data.json"); errors.Is(err, os.ErrNotExist) {
		f, err = os.Create(dataDir + "/data.json")
	} else {
		f, err = os.Open(dataDir + "/data.json")
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(&box.data)
	if err != nil && !errors.Is(err, io.EOF) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	box.filePath = dataDir + "/data.json"

	if box.data.NaoSet == nil {
		box.data.NaoSet = make(map[string]Set, 0)
	}

	return &box
}
