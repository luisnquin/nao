package store

import (
	"errors"
	"sync"
	"time"

	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/luisnquin/nao/v3/internal/store/keyutils"
	"github.com/luisnquin/nao/v3/internal/store/tagutils"
)

var (
	ErrNoAvailableNotes = errors.New("no available notes available")
	ErrNoteNotFound     = errors.New("note not found")
)

type NotesRepository struct {
	data *data.Buffer
	tag  tagutils.Tag
}

func NewNotesRepository(data *data.Buffer) NotesRepository {
	return NotesRepository{
		tag:  tagutils.New(data),
		data: data,
	}
}

func (r NotesRepository) LastAccessed() (models.Note, error) {
	note, ok := r.data.Notes[r.data.LastAccess]
	if !ok {
		return note, ErrNoteNotFound
	}

	note.Key = r.data.LastAccess

	return note, nil
}

func (r NotesRepository) Slice() []models.Note {
	notes := make([]models.Note, 0, len(r.data.Notes))

	// TODO: autorepair key

	for key, note := range r.data.Notes {
		note.Key = key

		notes = append(notes, note)
	}

	return notes
}

func (r NotesRepository) Iter() <-chan models.Note {
	ch := make(chan models.Note)
	mu := new(sync.RWMutex)

	go func() {
		mu.RLock()

		for key, note := range r.data.Notes {
			note.Key = key
			ch <- note
		}

		close(ch)
		mu.RUnlock()
	}()

	return ch
}

func (r NotesRepository) IterKey() <-chan string {
	ch := make(chan string)
	mu := new(sync.RWMutex)

	go func() {
		mu.RLock()

		for key := range r.data.Notes {
			ch <- key
		}

		close(ch)
		mu.RUnlock()
	}()

	return ch
}

func (r NotesRepository) ListAllKeys() []string {
	keys := make([]string, 0, len(r.data.Notes))

	for key := range r.data.Notes {
		keys = append(keys, key)
	}

	return keys
}

func (r NotesRepository) Total() int {
	return len(r.data.Notes)
}

func (r NotesRepository) Random() (models.Note, error) {
	for key, note := range r.data.Notes {
		note.Key = key

		return note, nil
	}

	return models.Note{}, ErrNoAvailableNotes
}

func (r NotesRepository) Get(key string) (models.Note, error) {
	note, ok := r.data.Notes[key]
	if !ok {
		return note, ErrNoteNotFound
	}

	r.data.LastAccess = key
	note.Key = key

	return note, r.data.Save()
}

func (r NotesRepository) TagExists(tag string) bool {
	for _, note := range r.data.Notes {
		if note.Tag == tag {
			return true
		}
	}

	return false
}

func (r NotesRepository) New(content, tag string, spentTime time.Duration) (string, error) {
	err := r.tag.IsValidAsNew(tag)
	if err != nil {
		return "", err
	}

	key := keyutils.New()

	r.data.Notes[key] = models.Note{
		Tag:        tag,
		Content:    content,
		CreatedAt:  time.Now(),
		LastUpdate: time.Now(),
		Version:    1,
	}

	return key, r.data.Save()
}

func (r NotesRepository) Replace(key string, note models.Note) error {
	_, ok := r.data.Notes[key]
	if !ok {
		return ErrNoteNotFound
	}

	r.data.Notes[key] = note

	return r.data.Save()
}

func (r NotesRepository) ModifyContent(key, content string, timeSpent time.Duration) error {
	note, ok := r.data.Notes[key]
	if !ok {
		return ErrNoteNotFound
	}

	note.LastUpdate = time.Now()
	note.TimeSpent += timeSpent
	note.Content = content
	note.Version++

	r.data.Notes[key] = note

	return r.data.Save()
}

func (r NotesRepository) ModifyTag(key, tag string) error {
	note, ok := r.data.Notes[key]
	if !ok {
		return ErrNoteNotFound
	}

	if err := r.tag.IsValidAsNew(tag); err != nil {
		return err
	}

	note.LastUpdate = time.Now()
	note.Tag = tag
	note.Version++

	r.data.Notes[key] = note

	return r.data.Save()
}

func (r NotesRepository) Delete(key string) error {
	_, ok := r.data.Notes[key]
	if !ok {
		return ErrNoteNotFound
	}

	delete(r.data.Notes, key)

	return r.data.Save()
}
