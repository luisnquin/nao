package data

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"
	"time"

	"github.com/cip8/autoname"
	"github.com/google/uuid"
)

var ErrSetNotFound error = errors.New("set not found")

func (d *Box) GetSet(key string) (Set, error) {
	set, ok := d.data.NaoSet[key]
	if !ok {
		return set, ErrSetNotFound
	}

	return set, nil
}

func (d *Box) ModifySet(key string, content string) error {
	set, ok := d.data.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	set.LastUpdate = time.Now()
	set.Content = content
	set.Version++

	d.data.NaoSet[key] = set

	return d.updateDataFile()
}

func (d *Box) ModifySetTag(key string, tag string) error {
	set, ok := d.data.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	set.LastUpdate = time.Now()
	set.Tag = tag
	set.Version++

	d.data.NaoSet[key] = set

	return d.updateDataFile()
}

func (d *Box) NewSet(content, contentType string) (string, error) {
	key := d.newKey()

	d.data.NaoSet[key] = Set{
		Tag:        autoname.Generate("-"),
		Content:    content,
		Type:       contentType,
		LastUpdate: time.Now(),
		Version:    1,
	}

	return key, d.updateDataFile()
}

func (d *Box) NewSetWithTag(content, contentType, tag string) (string, error) {
	key := d.newKey()

	d.data.NaoSet[key] = Set{
		Tag:        tag,
		Type:       contentType,
		Content:    content,
		LastUpdate: time.Now(),
		Version:    1,
	}

	return key, d.updateDataFile()
}

func (d *Box) SearchSetByKeyPattern(pattern string) (string, Set, error) {
	set, ok := d.data.NaoSet[pattern]
	if ok {
		return pattern, set, nil
	}

	for k, set := range d.data.NaoSet {
		if strings.HasPrefix(k, pattern) {
			return k, set, nil
		}
	}

	return "", set, ErrSetNotFound
}

func (d *Box) SearchSetByKeyTagPattern(pattern string) (string, Set, error) {
	set, ok := d.data.NaoSet[pattern]
	if ok {
		return pattern, set, nil
	}

	for k, set := range d.data.NaoSet {
		if strings.HasPrefix(k, pattern) {
			return k, set, nil
		}

		if strings.HasPrefix(set.Tag, pattern) {
			return k, set, nil
		}
	}

	return "", set, ErrSetNotFound
}

func (d *Box) DeleteSet(key string) error {
	_, ok := d.data.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	delete(d.data.NaoSet, key)

	return d.updateDataFile()
}

func (d *Box) GetMainSet() Set {
	return d.data.MainSet
}

func (d *Box) ModifyMainSet(content string) error {
	d.data.MainSet.Content = content
	d.data.MainSet.LastUpdate = time.Now()
	d.data.MainSet.Version++

	return d.updateDataFile()
}

func (d *Box) CleanMainSet() error {
	d.data.MainSet.Content = ""
	d.data.MainSet.LastUpdate = time.Now()
	d.data.MainSet.Version++

	return d.updateDataFile()
}

func (d *Box) ListSets() []SetView {
	sets := make([]SetView, 0)

	for k, v := range d.data.NaoSet {
		sets = append(sets, SetView{
			Key:        k,
			Tag:        v.Tag,
			Type:       v.Type,
			Content:    v.Content,
			Version:    v.Version,
			LastUpdate: v.LastUpdate,
		})
	}

	return sets
}

func (d *Box) ListSetWithHiddenContent() []SetViewWithoutContent {
	sets := make([]SetViewWithoutContent, 0)

	for k, v := range d.data.NaoSet {
		sets = append(sets, SetViewWithoutContent{
			Key:        k,
			Tag:        v.Tag,
			Type:       v.Type,
			Version:    v.Version,
			LastUpdate: v.LastUpdate,
		})
	}

	return sets
}

func (d *Box) ListAllKeys() []string {
	keys := make([]string, 0)

	for k := range d.data.NaoSet {
		keys = append(keys, k)
	}

	return keys
}

func (d *Box) updateDataFile() error {
	content, err := json.MarshalIndent(d.data, "", "\t")
	if err != nil {
		return err
	}

	/*
		if d.password != "" {
			content, err = security.EncryptContent([]byte(d.password), content)
			if err != nil {
				return err
			}
		}
	*/

	return ioutil.WriteFile(d.filePath, content, 0644)
}

func (d *Box) newKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
