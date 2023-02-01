package models

import (
	"time"

	"github.com/luisnquin/nao/v3/internal/utils"
)

type Note struct {
	Key        string        `json:"-"`
	Tag        string        `json:"tag,omitempty"`
	Content    string        `json:"content"`
	CreatedAt  time.Time     `json:"createdAt,omitempty"`
	LastUpdate time.Time     `json:"lastUpdate"`
	Version    int           `json:"version"`
	TimeSpent  time.Duration `json:"timeSpent"`
	// The number of get operations performed on a note.
	Picks uint64 `json:"picks"`
}

func (n *Note) Size() int {
	return utils.GetSize(n)
}

func (n *Note) ReadableSize() string {
	return utils.GetHumanReadableSize(n)
}
