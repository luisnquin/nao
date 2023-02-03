package note

import (
	"errors"
	"regexp"
	"strings"

	"github.com/luisnquin/nao/v3/internal/data"
)

type Tagger struct {
	data *data.Buffer
}

var rxTag = regexp.MustCompile(`^[A-z\_\-\@0-9]+$`)

var (
	ErrTagAlreadyExists = errors.New("tag already exists")
	ErrTagNotProvided   = errors.New("tag not provided")
	ErrTagNotFound      = errors.New("tag not found")
	ErrTagInvalid       = errors.New("tag invalid")
)

func NewTagger(data *data.Buffer) Tagger {
	return Tagger{data}
}

func (t Tagger) Like(tag string) (string, error) {
	for key, note := range t.data.Notes {
		if strings.HasPrefix(note.Tag, tag) {
			return key, nil
		}
	}

	return "", ErrTagNotFound
}

func (t Tagger) Exists(tag string) bool {
	for _, note := range t.data.Notes {
		if note.Tag == tag {
			return true
		}
	}

	return false
}

func (t Tagger) IsValid(tag string) error {
	if tag == "" {
		return ErrTagNotProvided
	}

	if !rxTag.MatchString(tag) {
		return ErrTagInvalid
	}

	return nil
}

func (t Tagger) IsValidAsNew(tag string) error {
	err := t.IsValid(tag)
	if err != nil {
		return err
	}

	if t.Exists(tag) {
		return ErrTagAlreadyExists
	}

	return nil
}
