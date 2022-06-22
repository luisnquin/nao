package data

import "time"

type (
	Box struct {
		data     BoxData
		password string
		filePath string
	}

	BoxData struct {
		LastAccess string         `json:"lastSet"`
		NaoSet     map[string]Set `json:"naoSet"`
		MainSet    Set            `json:"mainDraft"`
	}

	Set struct {
		Tag        string    `json:"tag,omitempty"`
		Content    string    `json:"content"`
		LastUpdate time.Time `json:"lastUpdate"`
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
		Content    string
		LastUpdate time.Time
	}

	SetViewWithoutContent struct {
		Tag        string
		Key        string
		LastUpdate time.Time
	}
)
