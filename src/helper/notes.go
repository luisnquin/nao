package helper

import "github.com/luisnquin/nao/src/data"

func SearchCriteriaInNoteView(notes []data.NoteView, f func(n data.NoteView) bool) bool {
	for _, note := range notes {
		if f(note) {
			return true
		}
	}

	return false
}
