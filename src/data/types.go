package data

import "time"

type (
	Box struct {
		box BoxData
	}

	BoxData struct {
		LastAccess string         `json:"lastSet"`
		NaoSet     map[string]Set `json:"naoSet"`
	}

	Set struct {
		Tag        string    `json:"tag,omitempty"`
		Type       string    `json:"type"`
		Content    string    `json:"content"`
		Extension  string    `json:"extension,omitempty"`
		Title      string    `json:"title,omitempty"`
		LastUpdate time.Time `json:"lastUpdate"`
		Version    int       `json:"version"`
	}
)

type Window struct {
	Hash       string
	Tag        string
	LastUpdate time.Time
}

type (
	SetView struct {
		Tag        string    `json:"tag"`
		Key        string    `json:"key"`
		Type       string    `json:"type"`
		Content    string    `json:"content"`
		Title      string    `json:"title"`
		Extension  string    `json:"extension"`
		LastUpdate time.Time `json:"lastUpdate"`
		Version    int       `json:"version"`
	}

	SetViewWithoutContent struct {
		Tag        string    `json:"tag"`
		Key        string    `json:"key"`
		Title      string    `json:"title"`
		Type       string    `json:"type"`
		Extension  string    `json:"extension"`
		LastUpdate time.Time `json:"lastUpdate"`
		Version    int       `json:"version"`
	}
)

type SetModifier interface {
	ModifySetContent(key string, content string) error
	ModifySetType(key string, sType string) error
}
