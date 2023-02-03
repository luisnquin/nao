package utils

import (
	"crypto/rand"
	"strings"

	"github.com/google/uuid"
)

func NewKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

// NewNanoID generates secure URL-friendly unique ID.
func NewNanoID() string {
	size := 20

	bytes := make([]byte, size)

	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}

	id := make([]rune, size)

	for i := 0; i < size; i++ {
		id[i] = []rune("..0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")[bytes[i]&61]
	}

	return string(id[:size])
}
