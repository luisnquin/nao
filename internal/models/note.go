package models

import (
	"time"

	"github.com/luisnquin/nao/internal/utils"
)

type Note struct {
	Key        string    `json:"-"`
	Tag        string    `json:"tag,omitempty"`
	Group      string    `json:"group"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content"`
	LastUpdate time.Time `json:"lastUpdate"`
	Version    int       `json:"version"`
}

func (n *Note) Size() int {
	return utils.GetSize(n)
}

func (n *Note) HumanReadableSize() string {
	return utils.GetHumanReadableSize(n)
}
