package data

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/goccy/go-json"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/models"
)

type Buffer struct {
	LastAccess string                 `json:"lastAccess,omitempty"`
	Notes      map[string]models.Note `json:"notes"`
	config     *config.AppConfig
}

func NewBuffer(config *config.AppConfig) (*Buffer, error) {
	data := Buffer{config: config}

	return &data, data.Load()
}

var (
	bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 0o5}
	key   = "ebee6254-d04e-4e51-be09-d0c7c8d4"
)

// Saves the current state of the data in the file. If the file
// doesn't exists then it will be created.
func (b *Buffer) Save() error {
	content, err := json.MarshalIndent(b, "", "\t")
	if err != nil {
		return fmt.Errorf("unexpected error, can't format data buffer to json: %w", err)
	}

	cipherText, err := encryptAndEncode(content, key)
	if err != nil {
		panic(err)
	}

	return ioutil.WriteFile(b.config.FS.DataFile, cipherText, 0o644)
}

// First data load, if there's no file to load then it creates it.
func (b *Buffer) Load() error {
	if err := b.Reload(); err != nil {
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

			return b.Reload()
		}

		return fmt.Errorf("unexpected error: %w", err)
	}

	return nil
}

// Reloads the data taking it from the expected file. If the file
// doesn't exists then throws an error and doesn't updates anything.
func (b *Buffer) Reload() error {
	file, err := os.Open(b.config.FS.DataFile)
	if err != nil {
		return err
	}

	// TODO: There should be something to migrate data.json to data.txt

	data, err := decryptAndDecode(file, key)
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

		if err = b.Save(); err != nil {
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
