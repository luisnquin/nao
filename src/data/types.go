package data

import "time"

type (
	Data struct {
		// Is the key of the set.
		LastAccess string         `json:"lastSet"`
		NaoSet     map[string]Set `json:"naoSet"`
		MainDraft  Set            `json:"mainDraft"`
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
