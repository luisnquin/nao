package security

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 0o5}

func DecryptAndDecode(stream io.Reader, key string) ([]byte, error) {
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

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(plainText, cipherText)

	return plainText, nil
}

func EncryptAndEncode(plainText []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(plainText))

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText, plainText)

	return []byte(base64.StdEncoding.EncodeToString(cipherText)), nil
}
