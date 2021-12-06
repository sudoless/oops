package oops_json

import (
	"encoding/json"
)

// Unmarshal calls json.Unmarshal but wraps the returning error in an oops.Error.
func Unmarshal(data []byte, v interface{}) error {
	return Error(json.Unmarshal(data, v))
}

// Marshal calls json.Marshal but wraps the returning error in an oops.Error.
func Marshal(v interface{}) ([]byte, error) {
	return ErrorM(json.Marshal(v))
}
