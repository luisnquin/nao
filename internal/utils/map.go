package utils

import (
	"bytes"
	"encoding/json"
)

func TypeToMapSlice(v any) []map[string]any {
	b := new(bytes.Buffer)
	r := make([]map[string]any, 0)

	_ = json.NewEncoder(b).Encode(v)
	_ = json.NewDecoder(b).Decode(&r)

	return r
}

func TypeToMap(v any) map[string]any {
	b := new(bytes.Buffer)
	r := make(map[string]any, 0)

	_ = json.NewEncoder(b).Encode(v)
	_ = json.NewDecoder(b).Decode(&r)

	return r
}
