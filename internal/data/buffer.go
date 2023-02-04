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
	"github.com/rs/zerolog"
	"github.com/zalando/go-keyring"
)

type (
	Buffer struct {
		Notes    map[string]models.Note `json:"notes"`
		Metadata Metadata               `json:"metadata"`
		log      *zerolog.Logger
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

func NewBuffer(logger *zerolog.Logger, config *config.Core) (*Buffer, error) {
	data := Buffer{config: config}

	if err := data.MigrateFileIfNeeded(); err != nil {
		return nil, err
	}

	return &data, data.Reload()
}

func (b *Buffer) MigrateFileIfNeeded() error {
	var srcFile, dstFile string

	b.log.Trace().Bool("encrypt", b.config.Encrypt).Send()

	if b.config.Encrypt {
		b.log.Trace().Msg("migration with projection from normal to encrypted")

		srcFile = b.config.FS.DataNormalFile
		dstFile = b.config.FS.DataEncryptedFile
	} else {
		b.log.Trace().Msg("migration with projection from encrypted to normal")

		srcFile = b.config.FS.DataEncryptedFile
		dstFile = b.config.FS.DataNormalFile
	}

	if !utils.FileExists(dstFile) && utils.FileExists(srcFile) {
		b.log.Trace().Str("source", srcFile).Str("destiny", dstFile).Msg("necessary migration, expected file not found")

		data, err := os.ReadFile(srcFile) // TODO: check that contains a valid json
		if err != nil {
			b.log.Err(err).Msg("unable to read source file")
			if os.IsPermission(err) {
				b.log.Error().Msg("it's a permissions error...")
			}

			return err
		}

		if b.config.Encrypt {
			secret := security.CreateRandomSecret()

			if err := security.SetSecretInKeyring(secret); err != nil {
				b.log.Err(err).Msg("failed attempt to set secret in keyring tool")

				return err
			}

			data, err = security.EncryptAndEncode(data, secret)
			if err != nil {
				b.log.Err(err).Msg("unable to encrypt and encode data by using SHA256 and base64, why?")

				return err
			}
		} else {
			secret, err := security.GetSecretFromKeyring()
			if err != nil {
				b.log.Err(err).
					Msg("failed attempt to get secret from keyring tool, this means that probably the data is irrecoverable")

				if errors.Is(err, keyring.ErrNotFound) {
					return errors.New("irrecoverable data file, secret not found")
				}

				return err
			}

			data, err = security.DecryptAndDecode(data, secret)
			if err != nil {
				b.log.Err(err).Msg("cannot decrypt the data file, maybe the file or secret is corrupted?")

				return err
			}

			b.log.Trace().Msg("deleting secret from keyring tool...")

			security.DeleteSecretFromKeyring()
		}

		b.log.Trace().Msg("data successfully recovered, the destiny file will be created and the other deleted")

		f, err := os.Create(dstFile)
		if err != nil {
			b.log.Err(err).Msg("failed attempt to create destiny file")
			if os.IsPermission(err) {
				b.log.Error().Msg("it's a permissions error...")
			}

			return err
		}

		b.log.Trace().Msg("deleting source file...")
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

	if keyToCare != "" && note.Version != 0 {
		if b.Notes == nil {
			b.Notes = make(map[string]models.Note, 1)
		}

		b.Notes[keyToCare] = note // TODO: suspect this
		b.Metadata = md
	}

	for k, n := range b.Notes { // TODO: log it
		if n.Tag == "" { // ? Or should I hide it in the ls command
			delete(b.Notes, k)
		}
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
