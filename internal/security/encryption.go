package security

import (
	"crypto/aes"
	"crypto/cipher"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 0o5}

// Decrypts the provided content by using AES-256 and also decodes
// the content using std base64.
func DecryptAndDecode(encryptedText []byte, key string) ([]byte, error) {
	encryptedText, err := DecodeFromBase64(encryptedText)
	if err != nil {
		return nil, err
	}

	return DecryptFromAES256(encryptedText, key)
}

// Encrypts the provided content by using AES-256 and also encodes
// the content using std base64.
func EncryptAndEncode(plainText []byte, key string) ([]byte, error) {
	encryptedText, err := EncryptToAES256(plainText, key)
	if err != nil {
		return nil, err
	}

	return EncodeToBase64(encryptedText), nil
}

func DecryptFromAES256(encryptedText []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	plainText := make([]byte, len(encryptedText))

	cipher.NewCFBDecrypter(block, iv).XORKeyStream(plainText, encryptedText)

	return plainText, nil
}

func EncryptToAES256(text []byte, key string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	cipherText := make([]byte, len(text))

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(cipherText, text)

	return cipherText, nil
}
