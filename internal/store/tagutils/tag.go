package tagutils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/luisnquin/nao/internal/data"
)

type Tag struct {
	data *data.Buffer
}

var (
	ErrTagAlreadyExists = errors.New("tag already exists")
	ErrTagNotProvided   = errors.New("tag not provided")
	ErrTagNotFound      = errors.New("tag not found")
	ErrTagInvalid       = errors.New("tag invalid")
)

func New(data *data.Buffer) Tag {
	return Tag{data}
}

func (t Tag) Like(tag string) (string, error) {
	for key, note := range t.data.Notes {
		if strings.HasPrefix(note.Tag, tag) {
			return key, nil
		}
	}

	return "", ErrTagNotFound
}

func (t Tag) Exists(tag string) bool {
	for _, note := range t.data.Notes {
		if note.Tag == tag {
			return true
		}
	}

	return false
}

func (t Tag) IsValid(tag string) error {
	if tag == "" {
		return ErrTagNotProvided
	}

	if !regexp.MustCompile(`^[A-z\_\-\@0-9]+$`).MatchString(tag) {
		return ErrTagInvalid
	}

	return nil
}

func (t Tag) IsValidAsNew(tag string) error {
	err := t.IsValid(tag)
	if err != nil {
		return err
	}

	if t.Exists(tag) {
		return ErrTagAlreadyExists
	}

	return nil
}
