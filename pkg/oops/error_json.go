package oops

import "encoding/json"

type errorJSON struct {
	Code    string   `json:"code"`
	Explain string   `json:"explain,omitempty"`
	Help    string   `json:"help,omitempty"`
	Multi   []string `json:"multi,omitempty"`
}

// MarshalJSON will encode the error in a format that is safe for a client/user to read without revealing any internal
// information about the structure or runtime of the program.
func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(&errorJSON{
		Code:    e.Error(),
		Explain: e.explanation.String(),
		Multi:   e.multi,
		Help:    e.help,
	})
}

// TODO, write unmarshal method
