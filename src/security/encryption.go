package security

import (
	"crypto/aes"
	"crypto/cipher"
	"math/rand"

	"golang.org/x/crypto/scrypt"
)

func EncryptContent(password, plainText []byte) ([]byte, error) {
	key, salt, err := NewKeyFromPassword(password, nil)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nonce, nonce, plainText, nil)
	cipherText = append(cipherText, salt...)

	return cipherText, nil
}

func DecryptContent(password, content []byte) ([]byte, error) {
	salt, content := content[len(content)-32:], content[:len(content)-32]

	key, _, err := NewKeyFromPassword(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := content[:gcm.NonceSize()], content[gcm.NonceSize():]

	return gcm.Open(nil, nonce, ciphertext, nil)
}

func NewKey() ([]byte, error) {
	key := make([]byte, 32)

	_, err := rand.Read(key)

	return key, err
}

func NewKeyFromPassword(password, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)

		_, err := rand.Read(salt)
		if err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key(password, salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
