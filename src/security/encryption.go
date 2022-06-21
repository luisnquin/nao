package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"math/rand"
)

func EncryptContent(password, content []byte) ([]byte, error) {
	key := NewKeyFromPassword(password)

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

	cipherText := gcm.Seal(nonce, nonce, content, nil)

	return cipherText, nil
}

func DecryptContent(password, content []byte) ([]byte, error) {
	key := NewKeyFromPassword(password)

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

func NewKeyFromPassword(password []byte) []byte {
	hasher := md5.New()
	hasher.Write(password)

	return []byte(hex.EncodeToString(hasher.Sum(nil)))
}
