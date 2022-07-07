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
	ErrTagNotProvided     error = errors.New("tag not provided")
	ErrInvalidSetType     error = errors.New("invalid set type")
	ErrGroupNotFound      error = errors.New("group not found")
	ErrKeyNotFound        error = errors.New("key not found")
	ErrSetNotFound        error = errors.New("set not found")
)

func (d *Box) Get(key string) (Note, error) {
	set, ok := d.box.NaoSet[key]
	if !ok {
		return set, ErrSetNotFound
	}

	d.box.LastAccess = key

	return set, d.updateFile()
}

func (d *Box) New(content, contentType string) (string, error) {
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

func (d *Box) NewWithTag(content, contentType, tag string) (string, error) {
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

func (d *Box) NewFrom(set Note) (string, error) {
	key := d.newKey()

	if set.Tag == "" {
		set.Tag = autoname.Generate("-")
	} else if d.TagAlreadyExists(set.Tag) {
		return "", ErrTagAlreadyExists
	}

	if set.Group != "" && !d.GroupExists(set.Group) {
		return "", ErrGroupNotFound
	}

	set.LastUpdate = time.Now()
	set.Version = 1

	if set.Type == constants.TypeMain && d.MainAlreadyExists() {
		return "", ErrMainAlreadyExists
	}

	d.box.NaoSet[key] = set

	return key, d.updateFile()
}

func (d *Box) ManyNewFrom(sets []Note) ([]string, error) {
	keys := make([]string, 0)

	for _, set := range sets {
		key, err := d.NewFrom(set)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (d *Box) Overwrite(key string, set Note) error {
	_, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	d.box.NaoSet[key] = set

	return d.updateFile()
}

func (d *Box) ReplaceBox(box BoxData) {
	d.box = box
}

func (d *Box) ModifyContent(key string, content string) error {
	set, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	set.LastUpdate = time.Now()
	set.History = append(set.History, Change{
		Key:       d.newKey(),
		Content:   set.Content,
		Timestamp: time.Now(),
	})
	set.Content = content
	set.Version++

	d.box.NaoSet[key] = set

	return d.updateFile()
}

func (d *Box) ModifyType(key string, sType string) error {
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

func (d *Box) ModifyTag(key string, tag string) error {
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

func (d *Box) Delete(key string) error {
	_, ok := d.box.NaoSet[key]
	if !ok {
		return ErrSetNotFound
	}

	delete(d.box.NaoSet, key)

	return d.updateFile()
}

func (d *Box) ResetToBefore(key string) error {
	s := d.box.NaoSet[key]

	if len(s.History) > 0 {
		s.Content = s.History[0].Content
		d.box.NaoSet[key] = s

		return d.updateFile()
	}

	return ErrKeyNotFound
}

func (d *Box) ResetTo(key, subKey string) error {
	s := d.box.NaoSet[key]

	for _, c := range s.History {
		if c.Key == subKey {
			s.Content = c.Content
			d.box.NaoSet[key] = s

			return d.updateFile()
		}
	}

	return ErrKeyNotFound
}

func (d *Box) ResetToWithDeletions(key, subKey string) error {
	var (
		s = d.box.NaoSet[key]
		t *time.Time
	)

	for _, c := range s.History {
		if c.Key == subKey {
			s.Content = c.Content
			d.box.NaoSet[key] = s

			t = &c.Timestamp

			break
		}
	}

	if t == nil {
		return ErrKeyNotFound
	}

	for i, c := range s.History {
		if c.Timestamp.After(*t) {
			s.History = append(s.History[:i], s.History[i+1:]...)
		}
	}

	d.box.NaoSet[key] = s

	return d.updateFile()
}

func (d *Box) SearchByKeyPattern(pattern string) (string, Note, error) {
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

func (d *Box) SearchByKeyTagPattern(pattern string) (string, Note, error) {
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

func (d *Box) List() []SetView {
	sets := make([]SetView, 0)

	for k, v := range d.box.NaoSet {
		sets = append(sets, SetView{
			LastUpdate: v.LastUpdate,
			Extension:  v.Extension,
			Version:    v.Version,
			Content:    v.Content,
			Group:      v.Group,
			Title:      v.Title,
			Type:       v.Type,
			Tag:        v.Tag,
			Key:        k,
		})
	}

	return sets
}

func (d *Box) ListWithHiddenContent() []SetViewWithoutContent {
	sets := make([]SetViewWithoutContent, 0)

	for k, v := range d.box.NaoSet {
		sets = append(sets, SetViewWithoutContent{
			LastUpdate: v.LastUpdate,
			Extension:  v.Extension,
			Version:    v.Version,
			Title:      v.Title,
			Group:      v.Group,
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
