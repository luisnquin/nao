package security

import "encoding/base64"

func DecodeFromBase64(encodedText []byte) ([]byte, error) {
	text := make([]byte, base64.StdEncoding.DecodedLen(len(encodedText)))
	n, err := base64.StdEncoding.Decode(text, encodedText)

	return text[:n], err
}

func EncodeToBase64(text []byte) []byte {
	encodedText := make([]byte, base64.StdEncoding.EncodedLen(len(text)))
	base64.StdEncoding.Encode(encodedText, text)

	return encodedText
}
