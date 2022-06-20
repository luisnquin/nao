package utils

import (
	"encoding/json"
	"io"
)

func EncodeToJSONIndent(b io.Writer, content any) error {
	e := json.NewEncoder(b)
	e.SetIndent("", "\t")
	return e.Encode(content)
}
