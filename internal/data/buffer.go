package data

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v3/internal"
	"github.com/luisnquin/nao/v3/internal/config"
	"github.com/luisnquin/nao/v3/internal/models"
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

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 0o5}

func NewBuffer(config *config.Core) (*Buffer, error) {
	data := Buffer{config: config}

	return &data, data.Reload()
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

	content, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return fmt.Errorf("unexpected error, can't format data buffer to json: %w", err)
	}

	caller, err := user.Current()
	if err != nil {
		panic(err)
	}

	key, err := keyring.Get(internal.AppName, caller.Username)
	if err != nil {
		return err
	}

	cipherText, err := encryptAndEncode(content, key)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(b.config.FS.DataFile, cipherText, 0o644)
}

// First data load, if there's no file to load then it creates it.
func (b *Buffer) Reload() error {
	if err := b.Load(); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(b.config.FS.DataDir, os.ModePerm)
			if err != nil {
				return fmt.Errorf("unable to create a new directory in '%s': %w", b.config.FS.DataFile, err)
			}

			file, err := os.Create(b.config.FS.DataFile)
			if err != nil {
				return fmt.Errorf("unable to create data file %s: %v", b.config.FS.DataFile, err)
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
	file, err := os.Open(b.config.FS.DataFile)
	if err != nil {
		return err
	}

	var data []byte

	if b.config.Encrypt {
		caller, err := user.Current()
		if err != nil {
			panic(err)
		}

		key, err := keyring.Get(internal.AppName, caller.Username)
		if err != nil {
			return err
		}

		data, err = decryptAndDecode(file, key)
	} else {
		data, err = io.ReadAll(file)
	}

	if err != nil {
		return err
	}

	defer file.Close()

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

func decryptAndDecode(stream io.Reader, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	encodedData, err := io.ReadAll(stream)
	if err != nil {
		return nil, err
	}

	cipherText, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(cipherText))

	cipher.NewCFBDecrypter(block, bytes).XORKeyStream(plainText, cipherText)

	return plainText, nil
}

func encryptAndEncode(plainText []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(plainText))

	cipher.NewCFBEncrypter(block, bytes).XORKeyStream(cipherText, plainText)

	return []byte(base64.StdEncoding.EncodeToString(cipherText)), nil
}

// Generates secure URL-friendly unique ID.
func generateRandomKey() []byte {
	size := 32

	bts := make([]byte, size)

	if _, err := rand.Read(bts); err != nil {
		panic(err)
	}

	id := make([]rune, size)

	for i := 0; i < size; i++ {
		id[i] = []rune("..0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")[bts[i]&61]
	}

	return []byte(string(id[:size]))
}
