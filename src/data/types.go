package data

import "time"

type (
	Box struct {
		box BoxData
	}

	BoxData struct {
		LastAccess string          `json:"lastSet"`
		NaoSet     map[string]Note `json:"naoSet"`
		Groups     []string        `json:"groups"`
	}

	Change struct {
		Key       string    `json:"key"`
		Content   string    `json:"content"`
		Timestamp time.Time `json:"timestamp"`
	}

	Note struct {
		Tag        string    `json:"tag,omitempty"`
		Type       string    `json:"type"`
		ShadowType string    `json:"shadowType"`
		Group      string    `json:"group"`
		Content    string    `json:"content"`
		Extension  string    `json:"extension,omitempty"`
		History    []Change  `json:"history"`
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

type NoteView struct {
	Tag        string    `json:"tag"`
	Key        string    `json:"key"`
	Type       string    `json:"type"`
	Size       string    `json:"size"`
	Content    string    `json:"content"`
	Group      string    `json:"group"`
	Title      string    `json:"title"`
	Extension  string    `json:"extension"`
	LastUpdate time.Time `json:"lastUpdate"`
	Version    int       `json:"version"`
}

type SetModifier interface {
	ModifyContent(key string, content string) error
	ModifyType(key string, sType string) error
}
