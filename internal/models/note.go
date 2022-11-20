package models

import (
	"time"

	"github.com/luisnquin/nao/v2/internal/utils"
)

type Note struct {
	Key        string    `json:"-"`
	Tag        string    `json:"tag,omitempty"`
	Title      string    `json:"title,omitempty"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	LastUpdate time.Time `json:"lastUpdate"`
	Version    int       `json:"version"`
}

func (n *Note) Size() int {
	return utils.GetSize(n)
}

func (n *Note) HumanReadableSize() string {
	return utils.GetHumanReadableSize(n)
}
