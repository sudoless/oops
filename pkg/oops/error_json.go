package oops

import "encoding/json"

type errorJSON struct {
	Code    string   `json:"code"`
	Explain string   `json:"explain,omitempty"`
	Help    string   `json:"help,omitempty"`
	Multi   []string `json:"multi,omitempty"`
}

func (e *Error) MarshalJSON() ([]byte, error) {
	return json.Marshal(&errorJSON{
		Code:    e.Error(),
		Explain: e.explanation.String(),
		Multi:   e.multi,
		Help:    e.help,
	})
}
