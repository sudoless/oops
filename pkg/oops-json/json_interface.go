package oops_json

import (
	"encoding/json"
)

// Unmarshal calls json.Unmarshal but wraps the returning error in an oops.Error.
func Unmarshal(data []byte, v any) error {
	return Error(json.Unmarshal(data, v))
}

// Marshal calls json.Marshal but wraps the returning error in an oops.Error.
func Marshal(v any) ([]byte, error) {
	return ErrorM(json.Marshal(v))
}
