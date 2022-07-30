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
	ErrMainAlreadyExists  error = errors.New("main note already exists")
	ErrGroupAlreadyExists error = errors.New("group already exists")
	ErrMainSetNotFound    error = errors.New("main note not found")
	ErrTagAlreadyExists   error = errors.New("tag already exists")
	ErrInvalidNoteType    error = errors.New("invalid note type")
	ErrTagNotProvided     error = errors.New("tag not provided")
	ErrGroupNotFound      error = errors.New("group not found")
	ErrNoteNotFound       error = errors.New("note not found")
	ErrKeyNotFound        error = errors.New("key not found")
	ErrTagInvalid         error = errors.New("tag invalid")
)

func (d *Box) Get(key string) (Note, error) {
	note, ok := d.box.NaoSet[key]
	if !ok {
		return note, ErrNoteNotFound
	}

	d.box.LastAccess = key

	return note, d.updateBoxFile()
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

	return key, d.updateBoxFile()
}

func (d *Box) NewWithTag(content, contentType, tag string) (string, error) {
	key := d.newKey()

	if tag == "" {
		return "", ErrTagNotProvided
	} else if err := d.TagIsValid(tag); err != nil {
		return "", err
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

	return key, d.updateBoxFile()
}

func (d *Box) GetLastKey() string {
	return d.box.LastAccess
}

func (d *Box) GetMainKey() (string, error) {
	for k, note := range d.box.NaoSet {
		if note.Type == constants.TypeMain {
			return k, nil
		}
	}

	return "", ErrMainSetNotFound
}

func (d *Box) NewFrom(note Note) (string, error) {
	key := d.newKey()

	if note.Tag == "" {
		note.Tag = autoname.Generate("-")
	} else if err := d.TagIsValid(note.Tag); err != nil {
		return "", err
	}

	if note.Group != "" && !d.GroupExists(note.Group) {
		return "", ErrGroupNotFound
	}

	note.LastUpdate = time.Now()
	note.Version = 1

	if note.Type == constants.TypeMain && d.MainAlreadyExists() {
		return "", ErrMainAlreadyExists
	}

	d.box.NaoSet[key] = note

	return key, d.updateBoxFile()
}

func (d *Box) ManyNewFrom(notes []Note) ([]string, error) {
	keys := make([]string, 0)

	for _, note := range notes {
		key, err := d.NewFrom(note)
		if err != nil {
			return nil, err
		}

		keys = append(keys, key)
	}

	return keys, nil
}

func (d *Box) Replace(key string, note Note) error {
	_, ok := d.box.NaoSet[key]
	if !ok {
		return ErrNoteNotFound
	}

	d.box.NaoSet[key] = note

	return d.updateBoxFile()
}

func (d *Box) ReplaceBox(box BoxData) {
	d.box = box
}

func (d *Box) ModifyContent(key string, content string) error {
	note, ok := d.box.NaoSet[key]
	if !ok {
		return ErrNoteNotFound
	}

	note.LastUpdate = time.Now()
	note.History = append(note.History, Change{
		Key:       d.newKey(),
		Content:   note.Content,
		Timestamp: time.Now(),
	})
	note.Content = content
	note.Version++

	d.box.NaoSet[key] = note

	return d.updateBoxFile()
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
		return ErrInvalidNoteType
	}

	note, ok := d.box.NaoSet[key]
	if !ok {
		return ErrNoteNotFound
	}

	note.Type = sType

	d.box.NaoSet[key] = note

	return d.updateBoxFile()
}

func (d *Box) ModifyTag(key string, tag string) error {
	note, ok := d.box.NaoSet[key]
	if !ok {
		return ErrNoteNotFound
	}

	if err := d.TagIsValid(tag); err != nil {
		return err
	}

	note.LastUpdate = time.Now()
	note.Tag = tag
	note.Version++

	d.box.NaoSet[key] = note

	return d.updateBoxFile()
}

func (d *Box) Delete(key string) error {
	_, ok := d.box.NaoSet[key]
	if !ok {
		return ErrNoteNotFound
	}

	delete(d.box.NaoSet, key)

	return d.updateBoxFile()
}

func (d *Box) ResetToBefore(key string) error {
	s := d.box.NaoSet[key]

	if len(s.History) > 0 {
		s.Content = s.History[0].Content
		d.box.NaoSet[key] = s

		return d.updateBoxFile()
	}

	return ErrKeyNotFound
}

func (d *Box) ResetTo(key, subKey string) error {
	s := d.box.NaoSet[key]

	for _, c := range s.History {
		if c.Key == subKey {
			s.Content = c.Content
			d.box.NaoSet[key] = s

			return d.updateBoxFile()
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

	return d.updateBoxFile()
}

func (d *Box) SearchByKeyPattern(pattern string) (string, Note, error) {
	note, ok := d.box.NaoSet[pattern]
	if ok {
		d.box.LastAccess = pattern
		return pattern, note, d.updateBoxFile()
	}

	for k, Note := range d.box.NaoSet {
		if strings.HasPrefix(k, pattern) {
			d.box.LastAccess = k
			return k, Note, d.updateBoxFile()
		}
	}

	return "", note, ErrNoteNotFound
}

func (d *Box) SearchByKeyTagPattern(pattern string) (string, Note, error) {
	note, ok := d.box.NaoSet[pattern]
	if ok {
		d.box.LastAccess = pattern
		return pattern, note, d.updateBoxFile()
	}

	for k, note := range d.box.NaoSet {
		if strings.HasPrefix(k, pattern) || strings.HasPrefix(note.Tag, pattern) {
			d.box.LastAccess = k
			return k, note, d.updateBoxFile()
		}
	}

	return "", note, ErrNoteNotFound
}

func (d *Box) List() []NoteView {
	notes := make([]NoteView, 0)

	for k, v := range d.box.NaoSet {
		notes = append(notes, NoteView{
			LastUpdate: v.LastUpdate,
			Extension:  v.Extension,
			Version:    v.Version,
			Content:    v.Content,
			Group:      v.Group,
			Title:      v.Title,
			Size:       d.boxSize(v),
			Type:       v.Type,
			Tag:        v.Tag,
			Key:        k,
		})
	}

	return notes
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
