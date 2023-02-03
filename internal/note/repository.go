package note

import (
	"sync"
	"time"

	"github.com/cip8/autoname"
	"github.com/luisnquin/nao/v3/internal/data"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/luisnquin/nao/v3/internal/utils"
)

type NotesRepository struct {
	data *data.Buffer
	tag  Tagger
}

func NewRepository(data *data.Buffer) NotesRepository {
	return NotesRepository{
		tag:  NewTagger(data),
		data: data,
	}
}

type ModifyOption func(*models.Note)

func WithTag(tag string) ModifyOption {
	return func(n *models.Note) {
		n.Tag = tag
	}
}

func WithSpentTime(duration time.Duration) ModifyOption {
	return func(n *models.Note) {
		n.TimeSpent += duration
	}
}

func WithContent(content string) ModifyOption {
	return func(n *models.Note) {
		n.Content = content
	}
}

func (r NotesRepository) Get(key string) (models.Note, error) {
	note, ok := r.data.Notes[key]
	if !ok {
		return note, ErrNoteNotFound
	}

	r.data.Metadata.LastAccess = data.KeyTag{
		Tag: note.Tag,
		Key: key,
	}

	note.Key = key

	return note, r.data.Save(key)
}

func (r NotesRepository) New(content string, modifiers ...ModifyOption) (string, error) {
	key := utils.NewKey()

	note := models.Note{
		Content:    content,
		CreatedAt:  time.Now(),
		LastUpdate: time.Now(),
		Version:    1,
	}

	for _, option := range modifiers {
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

func (r NotesRepository) Update(key string, modifiers ...ModifyOption) error {
	note, ok := r.data.Notes[key]
	if !ok {
		return ErrNoteNotFound
	}

	for _, option := range modifiers {
		option(&note)
	}

	note.LastUpdate = time.Now()
	note.Version++

	r.data.Notes[key] = note

	return r.data.Save(key)
}

func (r NotesRepository) Delete(key string) error {
	_, ok := r.data.Notes[key]
	if !ok {
		return ErrNoteNotFound
	}

	delete(r.data.Notes, key)

	return r.data.Save(key)
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

func (r NotesRepository) LastAccessed() (models.Note, error) {
	note, ok := r.data.Notes[r.data.Metadata.LastAccess.Key]
	if !ok {
		return note, ErrNoteNotFound
	}

	note.Key = r.data.Metadata.LastAccess.Key

	return note, nil
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

	return r.data.Save(key)
}

func (r NotesRepository) TagExists(tag string) bool {
	for _, note := range r.data.Notes {
		if note.Tag == tag {
			return true
		}
	}

	return false
}
