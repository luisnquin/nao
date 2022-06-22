package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewKey() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
