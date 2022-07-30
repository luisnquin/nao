package data

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/config"
	"github.com/luisnquin/nao/src/utils"
)

func (d *Box) updateBoxFile() error {
	content, err := json.MarshalIndent(d.box, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(config.App.Paths.DataFile, content, 0o644)
}

func (d *Box) newKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func (d *Box) TagIsValid(tag string) error {
	if tag != "" && !regexp.MustCompile(`^[A-z\_\-\@0-9]+$`).MatchString(tag) {
		return ErrTagInvalid
	}

	if ok := d.TagAlreadyExists(tag); ok {
		return ErrTagAlreadyExists
	}

	return nil
}

func (d *Box) boxSize(n Note) string {
	content, _ := json.Marshal(n)

	return utils.BytesToStorageUnits(int64(len(content)))
}
