package helper

import "github.com/luisnquin/nao/src/store"

func SearchCriteriaInNoteView(notes []store.NoteView, f func(n store.NoteView) bool) bool {
	for _, note := range notes {
		if f(note) {
			return true
		}
	}

	return false
}
