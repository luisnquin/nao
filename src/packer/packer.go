package packer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ProtonMail/go-appdir"
	"github.com/cip8/autoname"
	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/core"
	"github.com/luisnquin/nao/src/utils"
)

type (
	Data struct {
		NaoSet    map[string]Set `json:"naoSet"`
		MainDraft Set            `json:"mainDraft"`
	}

	Set struct {
		Name       string    `json:"name,omitempty"` // TODO: It should be named as 'Tag'
		Content    string    `json:"content"`
		LastUpdate time.Time `json:"lastUpdate"`
	}
)

type Window struct {
	Hash       string
	Name       string
	LastUpdate time.Time
}

func LoadMainDraft() (path string, remove func(), err error) {
	appDirs := appdir.New(core.AppName)

	data, err := LoadUserData(appDirs.UserData())
	if err != nil {
		return "", nil, err
	}

	file, remove := NewCachedIn(appDirs.UserCache())

	_, err = file.WriteString(data.MainDraft.Content)

	return file.Name(), remove, err
}

func OverwriteMainDraft(content []byte) error {
	appDirs := appdir.New(core.AppName)

	data, err := LoadUserData(appDirs.UserData())
	if err != nil {
		return err
	}

	data.MainDraft.Content = string(content)
	data.MainDraft.LastUpdate = time.Now()

	b := new(bytes.Buffer)

	err = utils.EncodeToJSONIndent(b, data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(appDirs.UserData()+"/data.json", b.Bytes(), 0644)
}

func NewCached() (f *os.File, remove func()) {
	cacheDir := appdir.New(core.AppName).UserCache()

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

func NewCachedIn(path string) (f *os.File, close func()) {
	_ = os.MkdirAll(path, os.ModePerm)

	file, err := os.Create(path + "/" + strings.ReplaceAll(uuid.NewString(), "-", "") + ".tmp")
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

func SaveContent(key string, content string) error { // TODO: add the capacibility to add names
	dataDir := appdir.New(core.AppName).UserData()

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
		Name:       autoname.Generate("-"),
		Content:    content,
		LastUpdate: time.Now(),
	}

	b := new(bytes.Buffer)

	err = utils.EncodeToJSONIndent(b, data)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(file.Name(), b.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func ListNaoSets() ([]Window, error) {
	list := make([]Window, 0)

	dataDir := appdir.New(core.AppName).UserData()

	if _, err := os.Stat(dataDir + "/data.json"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		return list, nil
	}

	file, err := os.Open(dataDir + "/data.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var data Data

	if err = json.NewDecoder(file).Decode(&data); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	for hash, set := range data.NaoSet {
		list = append(list, Window{
			Hash:       hash,
			Name:       set.Name,
			LastUpdate: set.LastUpdate,
		})
	}

	return list, nil
}

func ListRawSets() (map[string]Set, error) {
	dataDir := appdir.New(core.AppName).UserData()

	if _, err := os.Stat(dataDir + "/data.json"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}

		return make(map[string]Set), nil
	}

	file, err := os.Open(dataDir + "/data.json")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var data Data

	if err = json.NewDecoder(file).Decode(&data); err != nil {
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
	}

	return data.NaoSet, nil
}

func SearchInSet(prefix string) (string, Set, error) {
	sets, err := ListRawSets()
	if err != nil {
		return "", Set{}, err
	}

	for k, v := range sets {
		if strings.HasPrefix(k, prefix) {
			return k, v, nil
		}
	}

	return "", Set{}, fmt.Errorf("No such file: " + prefix)
}
