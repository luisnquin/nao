package store

import (
	"sync"
	"time"

	"github.com/cip8/autoname"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/luisnquin/nao/v3/internal/store/tagutils"
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
	note, ok := r.data.Notes[r.data.Metadata.LastAccess.Key]
	if !ok {
		return note, internal.ErrNoteNotFound
	}

	note.Key = r.data.Metadata.LastAccess.Key

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

func (r NotesRepository) Exists(key string) bool {
	_, ok := r.data.Notes[key]

	return ok
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

	return models.Note{}, internal.ErrNoAvailableNotes
}

func (r NotesRepository) Get(key string) (models.Note, error) {
	note, ok := r.data.Notes[key]
	if !ok {
		return note, internal.ErrNoteNotFound
	}

	r.data.Metadata.LastAccess = data.KeyTag{
		Tag: note.Tag,
		Key: key,
	}

	note.Key = key

	return note, r.data.Save(key)
}

func (r NotesRepository) TagExists(tag string) bool {
	for _, note := range r.data.Notes {
		if note.Tag == tag {
			return true
		}
	}

	return false
}

type Option func(*models.Note)

func WithTag(tag string) Option {
	return func(n *models.Note) {
		n.Tag = tag
	}
}

func WithSpentTime(duration time.Duration) Option {
	return func(n *models.Note) {
		n.TimeSpent += duration
	}
}

func (r NotesRepository) New(content string, options ...Option) (string, error) {
	key := internal.NewKey()

	note := models.Note{
		Content:    content,
		CreatedAt:  time.Now(),
		LastUpdate: time.Now(),
		Version:    1,
	}

	for _, option := range options {
		option(&note)
	}

	if note.Tag == "" {
		note.Tag = autoname.Generate("-")
	} else {
		if err := r.tag.IsValidAsNew(note.Tag); err != nil {
			return "", err
		}
	}

	r.data.Metadata.LastCreated = data.KeyTag{
		Tag: note.Tag,
		Key: key,
	}

	r.data.Notes[key] = note

	return key, r.data.Save(key)
}

func (r NotesRepository) Replace(key string, note models.Note) error {
	_, ok := r.data.Notes[key]
	if !ok {
		return internal.ErrNoteNotFound
	}

	r.data.Notes[key] = note

	return r.data.Save(key)
}

func (r NotesRepository) ModifyContent(key, content string, options ...Option) error {
	note, ok := r.data.Notes[key]
	if !ok {
		return internal.ErrNoteNotFound
	}

	for _, option := range options {
		option(&note)
	}

	note.LastUpdate = time.Now()
	note.Content = content
	note.Version++

	r.data.Notes[key] = note

	return r.data.Save(key)
}

func (r NotesRepository) ModifyTag(key, tag string) error {
	note, ok := r.data.Notes[key]
	if !ok {
		return internal.ErrNoteNotFound
	}

	if err := r.tag.IsValidAsNew(tag); err != nil {
		return err
	}

	note.LastUpdate = time.Now()
	note.Tag = tag
	note.Version++

	r.data.Notes[key] = note

	return r.data.Save(key)
}

func (r NotesRepository) Delete(key string) error {
	_, ok := r.data.Notes[key]
	if !ok {
		return internal.ErrNoteNotFound
	}

	delete(r.data.Notes, key)

	return r.data.Save(key)
}
