package data

import (
	"errors"
	"strings"
	"time"

	"github.com/cip8/autoname"
	"github.com/luisnquin/nao/src/constants"
	"github.com/luisnquin/nao/src/utils"
)

var (
	ErrMainAlreadyExists  error = errors.New("main set already exists")
	ErrGroupAlreadyExists error = errors.New("group already exists")
	ErrMainSetNotFound    error = errors.New("main set not found")
	ErrTagAlreadyExists   error = errors.New("tag already exists")
	ErrGroupNotFound      error = errors.New("group not found")
	ErrTagNotProvided     error = errors.New("tag not provided")
	ErrInvalidSetType     error = errors.New("invalid set type")
	ErrSetNotFound        error = errors.New("set not found")
)

func (d *Box) ModifyBox(box BoxData) {
	d.box = box
}

func (d *Box) GetSet(key string) (Note, error) {
	set, ok := d.box.NaoSet[key]
	if !ok {
		return set, ErrSetNotFound
	}

	d.box.LastAccess = key

	return set, d.updateFile()
}

func (d *Box) GetGroups() []string {
	return d.box.Groups
}

func (d *Box) GetLastKey() string {
	return d.box.LastAccess
}

func (d *Box) GetMainKey() (string, error) {
	for k, set := range d.box.NaoSet {
		if set.Type == constants.TypeMain {
			return k, nil
		}
	}

	return "", ErrMainSetNotFound
}

func (d *Box) NewGroup(name string) error {
	for _, group := range d.box.Groups {
		if group == name {
			return ErrGroupAlreadyExists
		}
	}

	d.box.Groups = append(d.box.Groups, name)

	return d.updateFile()
}

func (d *Box) NewSet(content, contentType string) (string, error) {
	key := d.newKey()

	if contentType == constants.TypeMain && d.MainAlreadyExists() {
		return "", ErrMainAlreadyExists
	}

	d.box.NaoSet[key] = Note{
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

	d.box.NaoSet[key] = Note{
		Tag:        tag,
		Type:       contentType,
		Content:    content,
		LastUpdate: time.Now(),
		Version:    1,
	}

	return key, d.updateFile()
}

func (d *Box) NewFromSet(set Note) (string, error) {
	key := d.newKey()

	if set.Tag == "" {
		set.Tag = autoname.Generate("-")
	} else if d.TagAlreadyExists(set.Tag) {
		return "", ErrTagAlreadyExists
	}

	set.LastUpdate = time.Now()
	set.Version = 1

	if set.Type == constants.TypeMain && d.MainAlreadyExists() {
		return "", ErrMainAlreadyExists
	}

	d.box.NaoSet[key] = set

	return key, d.updateFile()
}

func (d *Box) NewSetsFromOutside(sets []Note) ([]string, error) {
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

func (d *Box) OverwriteSet(key string, set Note) error {
	_, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	d.box.NaoSet[key] = set

	return d.updateFile()
}

func (d *Box) ModifySetContent(key string, content string) error {
	set, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	set.LastUpdate = time.Now()
	set.Content = content
	set.Version++

	d.box.NaoSet[key] = set

	return d.updateFile()
}

func (d *Box) ModifySetType(key string, sType string) error {
	validTypes := []string{
		constants.TypeDefault,
		constants.TypeImported,
		constants.TypeMain,
		constants.TypeMerged,
	}

	if sType == constants.TypeMain && d.MainAlreadyExists() {
		return ErrMainAlreadyExists
	}

	if !utils.Contains(validTypes, sType) {
		return ErrInvalidSetType
	}

	set, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	set.Type = sType

	d.box.NaoSet[key] = set

	return d.updateFile()
}

func (d *Box) ModifySetTag(key string, tag string) error {
	set, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	if d.TagAlreadyExists(tag) {
		return ErrTagAlreadyExists
	}

	set.LastUpdate = time.Now()
	set.Tag = tag
	set.Version++

	d.box.NaoSet[key] = set

	return d.updateFile()
}

func (d *Box) ModifyGroupName(oldName, newName string) error {
	for i, group := range d.box.Groups {
		if group == oldName {
			d.box.Groups[i] = newName

			for k, set := range d.box.NaoSet {
				if set.Group == oldName {
					set.Group = newName

					d.box.NaoSet[k] = set
				}
			}

			return d.updateFile()
		}
	}

	return ErrGroupNotFound
}

func (d *Box) DeleteGroupWithRelated(name string) error {
	for i, group := range d.box.Groups {
		if group == name {
			d.box.Groups = append(d.box.Groups[:i], d.box.Groups[i+1:]...)

			for k, set := range d.box.NaoSet {
				if set.Group == name {
					delete(d.box.NaoSet, k)
				}
			}

			return d.updateFile()
		}
	}

	return ErrGroupNotFound
}

func (d *Box) DeleteGroup(name string) error {
	for i, group := range d.box.Groups {
		if group == name {
			d.box.Groups = append(d.box.Groups[:i], d.box.Groups[i+1:]...)

			for k, set := range d.box.NaoSet {
				if set.Group == name {
					set.Group = ""
					d.box.NaoSet[k] = set
				}
			}

			return d.updateFile()
		}
	}

	return ErrGroupNotFound
}

func (d *Box) DeleteSet(key string) error {
	_, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	delete(d.box.NaoSet, key)

	return d.updateFile()
}

func (d *Box) SearchSetByKeyPattern(pattern string) (string, Note, error) {
	set, ok := d.box.NaoSet[pattern]
	if ok {
		d.box.LastAccess = pattern
		return pattern, set, d.updateFile()
	}

	for k, set := range d.box.NaoSet {
		if strings.HasPrefix(k, pattern) {
			d.box.LastAccess = k
			return k, set, d.updateFile()
		}
	}

	return "", set, ErrSetNotFound
}

func (d *Box) SearchSetByKeyTagPattern(pattern string) (string, Note, error) {
	set, ok := d.box.NaoSet[pattern]
	if ok {
		d.box.LastAccess = pattern
		return pattern, set, d.updateFile()
	}

	for k, set := range d.box.NaoSet {
		if strings.HasPrefix(k, pattern) || strings.HasPrefix(set.Tag, pattern) {
			d.box.LastAccess = k
			return k, set, d.updateFile()
		}
	}

	return "", set, ErrSetNotFound
}

func (d *Box) ListSets() []SetView {
	sets := make([]SetView, 0)

	for k, v := range d.box.NaoSet {
		sets = append(sets, SetView{
			LastUpdate: v.LastUpdate,
			Extension:  v.Extension,
			Version:    v.Version,
			Content:    v.Content,
			Title:      v.Title,
			Type:       v.Type,
			Tag:        v.Tag,
			Key:        k,
		})
	}

	return sets
}

func (d *Box) ListSetWithHiddenContent() []SetViewWithoutContent {
	sets := make([]SetViewWithoutContent, 0)

	for k, v := range d.box.NaoSet {
		sets = append(sets, SetViewWithoutContent{
			LastUpdate: v.LastUpdate,
			Extension:  v.Extension,
			Version:    v.Version,
			Title:      v.Title,
			Type:       v.Type,
			Tag:        v.Tag,
			Key:        k,
		})
	}

	return sets
}

func (d *Box) ListAllKeys() []string {
	keys := make([]string, 0)

	for k := range d.box.NaoSet {
		keys = append(keys, k)
	}

	return keys
}

func (d *Box) TagAlreadyExists(tag string) bool {
	for _, s := range d.box.NaoSet {
		if s.Tag == tag {
			return true
		}
	}

	return false
}

func (d *Box) MainAlreadyExists() bool {
	for _, s := range d.box.NaoSet {
		if s.Type == constants.TypeMain {
			return true
		}
	}

	return false
}
