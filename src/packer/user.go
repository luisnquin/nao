package packer

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

func LoadUserData(dirPath string) (Data, error) {
	var data Data

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return data, err
	}

	var file *os.File

	if _, err = os.Stat(dirPath + "/data.json"); errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(dirPath + "/data.json")
	} else {
		file, err = os.Open(dirPath + "/data.json")
	}

	if err != nil {
		return data, err
	}

	err = json.NewDecoder(file).Decode(&data)
	if err != nil && !errors.Is(err, io.EOF) {
		return data, err
	}

	err = file.Close()

	return data, err
}
