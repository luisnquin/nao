package packer

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ProtonMail/go-appdir"
	"github.com/google/uuid"
)

type Data struct {
	NaoSet map[string]Set `json:"naoSet"`
}

type Set struct {
	LastUpdate time.Time `json:"lastUpdate"`
	Content    string    `json:"content"`
}

func NewCached() (f *os.File, close func()) {
	cacheDir := appdir.New("nao").UserCache()

	err := os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(cacheDir + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
	if err != nil {
		panic(err)
	}

	return file, func() {
		err = file.Close()
		if err != nil {
			panic(err)
		}

		err = os.Remove(file.Name())
		if err != nil {
			panic(err)
		}
	}
}

func SaveContent(key string, content string) error {
	dataDir := appdir.New("nao").UserData()

	err := os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return err
	}

	var file *os.File

	if _, err := os.Stat(dataDir + "/data.json"); errors.Is(err, os.ErrNotExist) {
		file, err = os.Create(dataDir + "/data.json")
	} else {
		file, err = os.Open(dataDir + "/data.json")
	}

	if err != nil {
		return err
	}

	defer file.Close()

	var data Data

	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		if !errors.Is(err, io.EOF) {
			return err
		}

		data.NaoSet = make(map[string]Set)
	}

	data.NaoSet[key] = Set{
		LastUpdate: time.Now(),
		Content:    content,
	}

	b := new(bytes.Buffer)

	e := json.NewEncoder(b)
	e.SetIndent("", "\t")

	err = e.Encode(data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file.Name(), b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
