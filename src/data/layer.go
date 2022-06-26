package data

import (
	"errors"
	"strings"
	"time"

	"github.com/cip8/autoname"
	"github.com/luisnquin/nao/src/constants"
)

var (
	ErrMainAlreadyExists error = errors.New("main set already exists")
	ErrMainSetNotFound   error = errors.New("main set not found")
	ErrTagAlreadyExists  error = errors.New("tag already exists")
	ErrTagNotProvided    error = errors.New("tag not provided")
	ErrSetNotFound       error = errors.New("set not found")
)

func (d *Box) GetSet(key string) (Set, error) {
	set, ok := d.data.NaoSet[key]
	if !ok {
		return set, ErrSetNotFound
	}

	d.data.LastAccess = key

	return set, d.updateFile()
}

func (d *Box) GetLastKey() string {
	return d.data.LastAccess
}

func (d *Box) GetMainKey() (string, error) {
	for k, set := range d.data.NaoSet {
		if set.Type == constants.TypeMain {
			return k, nil
		}
	}

	return "", ErrMainSetNotFound
}

func (d *Box) NewSet(content, contentType string) (string, error) {
	key := d.newKey()

	if contentType == constants.TypeMain && d.MainAlreadyExists() {
		return "", ErrMainAlreadyExists
	}

	d.data.NaoSet[key] = Set{
		Tag:        autoname.Generate("-"),
		Content:    content,
		Type:       contentType,
		LastUpdate: time.Now(),
		Version:    1,
	}

	return key, d.updateFile()
}

func (d *Box) NewSetWithTag(content, contentType, tag string) (string, error) {
	key := d.newKey()

	if tag == "" {
		return "", ErrTagNotProvided
	}

	if d.TagAlreadyExists(tag) {
		return "", ErrTagAlreadyExists
	}

	if contentType == constants.TypeMain && d.MainAlreadyExists() {
		return "", ErrMainAlreadyExists
	}

	d.data.NaoSet[key] = Set{
		Tag:        tag,
		Type:       contentType,
		Content:    content,
		LastUpdate: time.Now(),
		Version:    1,
	}

	return key, d.updateFile()
}

func (d *Box) NewFromSet(set Set) (string, error) {
	key := d.newKey()

	if set.Tag == "" {
		set.Tag = autoname.Generate("-")
	}

	set.LastUpdate = time.Now()
	set.Version = 1

	if set.Type == constants.TypeMain && d.MainAlreadyExists() {
		return "", ErrMainAlreadyExists
	}

	d.data.NaoSet[key] = set

	return key, d.updateFile()
}

func (d *Box) NewSetsFromOutside(sets []Set) ([]string, error) {
	keys := make([]string, 0)

	for _, set := range sets {
		key, err := d.NewFromSet(set)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
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

	return d.updateFile()
}

func (d *Box) ModifySetTag(key string, tag string) error {
	set, ok := d.data.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	if d.TagAlreadyExists(tag) {
		return ErrTagAlreadyExists
	}

	set.LastUpdate = time.Now()
	set.Tag = tag
	set.Version++

	d.data.NaoSet[key] = set

	return d.updateFile()
}

func (d *Box) DeleteSet(key string) error {
	_, ok := d.data.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	delete(d.data.NaoSet, key)

	return d.updateFile()
}

func (d *Box) SearchSetByKeyPattern(pattern string) (string, Set, error) {
	set, ok := d.data.NaoSet[pattern]
	if ok {
		d.data.LastAccess = pattern
		return pattern, set, d.updateFile()
	}

	for k, set := range d.data.NaoSet {
		if strings.HasPrefix(k, pattern) {
			d.data.LastAccess = k
			return k, set, d.updateFile()
		}
	}

	return "", set, ErrSetNotFound
}

func (d *Box) SearchSetByKeyTagPattern(pattern string) (string, Set, error) {
	set, ok := d.data.NaoSet[pattern]
	if ok {
		d.data.LastAccess = pattern
		return pattern, set, d.updateFile()
	}

	for k, set := range d.data.NaoSet {
		if strings.HasPrefix(k, pattern) || strings.HasPrefix(set.Tag, pattern) {
			d.data.LastAccess = k
			return k, set, d.updateFile()
		}
	}

	return "", set, ErrSetNotFound
}

func (d *Box) ListSets() []SetView {
	sets := make([]SetView, 0)

	for k, v := range d.data.NaoSet {
		sets = append(sets, SetView{
			LastUpdate: v.LastUpdate,
			Extension:  v.Extension,
			Version:    v.Version,
			Content:    v.Content,
			Type:       v.Type,
			Tag:        v.Tag,
			Key:        k,
		})
	}

	return sets
}

func (d *Box) ListSetWithHiddenContent() []SetViewWithoutContent {
	sets := make([]SetViewWithoutContent, 0)

	for k, v := range d.data.NaoSet {
		sets = append(sets, SetViewWithoutContent{
			LastUpdate: v.LastUpdate,
			Extension:  v.Extension,
			Version:    v.Version,
			Type:       v.Type,
			Tag:        v.Tag,
			Key:        k,
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

func (d *Box) TagAlreadyExists(tag string) bool {
	for _, s := range d.data.NaoSet {
		if s.Tag == tag {
			return true
		}
	}

	return false
}

func (d *Box) MainAlreadyExists() bool {
	for _, s := range d.data.NaoSet {
		if s.Type == constants.TypeMain {
			return true
		}
	}

	return false
}
