package data

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/config"
)

func (d *Box) updateFile() error {
	content, err := json.MarshalIndent(d.box, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(config.App.Paths.DataFile, content, 0644)
}

func (d *Box) newKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
