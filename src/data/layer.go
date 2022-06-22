package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cip8/autoname"
	"github.com/google/uuid"
	"github.com/luisnquin/nao/src/security"
)

var ErrSetNotFound error = errors.New("set not found")

func (d *Box) GetSet(key string) (Set, error) {
	set, ok := d.NaoSet[key]
	if !ok {
		return set, ErrSetNotFound
	}

	return set, nil
}

func (d *Box) ModifySet(key string, content string) error {
	set, ok := d.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	set.LastUpdate = time.Now()
	set.Content = content

	d.NaoSet[key] = set

	d.updateDataFile()

	return nil
}

func (d *Box) ModifySetTag(key string, tag string) error {
	set, ok := d.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	set.Tag = tag

	d.NaoSet[key] = set

	d.updateDataFile()

	return nil
}

func (d *Box) NewSet(content string) (string, error) {
	key := d.newKey()

	d.NaoSet[key] = Set{
		Tag:        autoname.Generate("-"),
		Content:    content,
		LastUpdate: time.Now(),
	}

	d.updateDataFile()

	return key, nil
}

func (d *Box) NewSetWithTag(content, tag string) (string, error) {
	key := d.newKey()

	d.NaoSet[key] = Set{
		Tag:        tag,
		Content:    content,
		LastUpdate: time.Now(),
	}

	d.updateDataFile()

	return key, nil
}

func (d *Box) SearchSetByPattern(pattern string) (string, Set, error) {
	set, ok := d.NaoSet[pattern]
	if ok {
		return pattern, set, nil
	}

	for k, set := range d.NaoSet {
		if strings.HasPrefix(k, pattern) {
			return k, set, nil
		}
	}

	return "", set, ErrSetNotFound
}

func (d *Box) ListSets() []SetView {
	sets := make([]SetView, 0)

	for k, v := range d.NaoSet {
		sets = append(sets, SetView{
			Key:        k,
			Tag:        v.Tag,
			Content:    v.Content,
			LastUpdate: v.LastUpdate,
		})
	}

	return sets
}

func (d *Box) ListSetWithHiddenContent() []SetViewWithoutContent {
	sets := make([]SetViewWithoutContent, 0)

	for k, v := range d.NaoSet {
		sets = append(sets, SetViewWithoutContent{
			Key:        k,
			Tag:        v.Tag,
			LastUpdate: v.LastUpdate,
		})
	}

	return sets
}

func (d *Box) GetMainSet() Set {
	return d.MainSet
}

func (d *Box) ModifyMainNote(content string) error {
	d.MainSet.Content = content

	d.updateDataFile()

	return nil
}

func (d *Box) updateDataFile() error {
	content, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}

	if d.password != "" {
		content, err = security.EncryptContent([]byte(d.password), content)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(d.filePath, content, 0644)
}

func (d *Box) newKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
