package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
)

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
	// iv is always stored in the encrypted text
	iv := encryptedText[:aes.BlockSize]
	encryptedText = encryptedText[aes.BlockSize:]

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	if len(encryptedText) < aes.BlockSize { // ? maybe remove/improve this
		return nil, errors.New("ciphertext too short")
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

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	encryptedText := make([]byte, len(text))

	cipher.NewCFBEncrypter(block, iv).XORKeyStream(encryptedText, text)

	// concatenate the iv and the ciphertext
	cipherWithVi := make([]byte, aes.BlockSize+len(encryptedText))
	copy(cipherWithVi[:aes.BlockSize], iv)
	copy(cipherWithVi[aes.BlockSize:], encryptedText)

	return cipherWithVi, nil
}
