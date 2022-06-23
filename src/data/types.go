package data

import "time"

type (
	Box struct {
		data     BoxData
		password string
	}

	BoxData struct {
		LastAccess string         `json:"lastSet"`
		NaoSet     map[string]Set `json:"naoSet"`
		MainSet    Set            `json:"mainDraft"`
	}

	Set struct {
		Tag        string    `json:"tag,omitempty"`
		Type       string    `json:"type"`
		Content    string    `json:"content"`
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
		Tag        string
		Key        string
		Type       string
		Content    string
		LastUpdate time.Time
		Version    int
	}

	SetViewWithoutContent struct {
		Tag        string
		Key        string
		Type       string
		LastUpdate time.Time
		Version    int
	}
)
