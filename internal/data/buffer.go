package data

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/models"
	"github.com/luisnquin/nao/v3/internal/security"
	"github.com/luisnquin/nao/v3/internal/utils"
	"github.com/zalando/go-keyring"
)

type (
	Buffer struct {
		Notes    map[string]models.Note `json:"notes"`
		Metadata Metadata               `json:"metadata"`
		config   *config.Core
	}

	Metadata struct {
		// The key of the last accessed note.
		LastCreated KeyTag `json:"lastCreated,omitempty"`
		// The key of the last accessed/modified note.
		LastAccess KeyTag `json:"lastAccess,omitempty"`
	}

	KeyTag struct {
		Key string `json:"key"`
		Tag string `json:"tag"`
	}
)

func NewBuffer(config *config.Core) (*Buffer, error) {
	data := Buffer{config: config}

	if err := data.MigrateFileIfNeeded(); err != nil {
		return nil, err
	}

	return &data, data.Reload()
}

func (b *Buffer) MigrateFileIfNeeded() error {
	var srcFile, dstFile string

	if b.config.Encrypt {
		srcFile = b.config.FS.DataNormalFile
		dstFile = b.config.FS.DataEncryptedFile
	} else {
		srcFile = b.config.FS.DataEncryptedFile
		dstFile = b.config.FS.DataNormalFile
	}

	if !utils.FileExists(dstFile) && utils.FileExists(srcFile) {
		data, err := os.ReadFile(srcFile) // TODO: check that contains a valid json
		if err != nil {
			return err
		}

		if b.config.Encrypt {
			secret := security.CreateRandomSecret()

			if err := security.SetSecretInKeyring(secret); err != nil {
				return err
			}

			data, err = security.EncryptAndEncode(data, secret)
			if err != nil {
				return err
			}
		} else {
			secret, err := security.GetSecretFromKeyring()
			if err != nil {
				if errors.Is(err, keyring.ErrNotFound) {
					return errors.New("irrecoverable data file, secret not found")
				}

				return err
			}

			data, err = security.DecryptAndDecode(data, secret)
			if err != nil {
				return err
			}

			security.DeleteSecretFromKeyring()
		}

		f, err := os.Create(dstFile)
		if err != nil {
			return err
		}

		os.Remove(srcFile)
		f.Write(data)
		f.Close()
	}

	return nil
}

// Saves the current state of the data in the file. If the file
// doesn't exists then it will be created.
func (b *Buffer) Commit(keyToCare string) error {
	note := b.Notes[keyToCare]
	md := b.Metadata

	if err := b.Reload(); err != nil {
		return err
	}

	if keyToCare != "" {
		if b.Notes == nil {
			b.Notes = make(map[string]models.Note, 1)
		}

		b.Notes[keyToCare] = note // TODO: suspect this
		b.Metadata = md
	}

	data, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return fmt.Errorf("unexpected error, can't format data buffer to json: %w", err)
	}

	if b.config.Encrypt {
		secret, err := security.GetSecretFromKeyring()
		if err != nil {
			return err
		}

		data, err = security.EncryptAndEncode(data, secret)
		if err != nil {
			return err
		}
	}

	return ioutil.WriteFile(b.config.FS.DataFile(b.config.Encrypt), data, internal.PermReadWrite)
}

// First data load, if there's no file to load then it creates it.
func (b *Buffer) Reload() error {
	if err := b.Load(); err != nil {
		if os.IsNotExist(err) {
			dataFile := b.config.FS.DataFile(b.config.Encrypt)

			err = os.MkdirAll(b.config.FS.DataDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("unable to create a new directory in '%s': %w", dataFile, err)
			}

			file, err := os.Create(dataFile)
			if err != nil {
				return fmt.Errorf("unable to create data file %s: %v", dataFile, err)
			}

			if !b.config.Encrypt {
				file.WriteString("{}") // Of course an empty file will never be a valid JSON
			}

			err = file.Close()
			if err != nil {
				return err
			}

			return b.Load()
		}

		return fmt.Errorf("unexpected error: %w", err)
	}

	return nil
}

// Reloads the data taking it from the expected file. If the file
// doesn't exists then throws an error and doesn't updates anything.
func (b *Buffer) Load() error {
	data, err := os.ReadFile(b.config.FS.DataFile(b.config.Encrypt))
	if err != nil {
		return err
	}

	if b.config.Encrypt {
		secret, err := security.GetSecretFromKeyring()
		if err != nil {
			return err
		}

		data, err = security.DecryptAndDecode(data, secret)
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(data, b)
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("unreadable json file: %w", err)
	}

	if b.Notes == nil {
		b.Notes = make(map[string]models.Note)

		if err = b.Commit(""); err != nil {
			return err
		}
	}

	return nil
}
