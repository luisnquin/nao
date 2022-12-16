package keyutils

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/luisnquin/nao/v3/internal/data"
)

type Key struct {
	data *data.Buffer
}

var ErrKeyNotFound = errors.New("key not found")

func NewDispatcher(data *data.Buffer) Key {
	return Key{data: data}
}

func New() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

func (k Key) Like(pattern string) (string, error) {
	_, ok := k.data.Notes[pattern]
	if ok {
		k.data.LastAccess = pattern

		return pattern, k.data.Save()
	}

	for key := range k.data.Notes {
		if strings.HasPrefix(key, pattern) {
			k.data.LastAccess = key

			return key, k.data.Save()
		}
	}

	return "", ErrKeyNotFound
}
