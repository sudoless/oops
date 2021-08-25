package oops

import (
	"encoding/json"
	"strings"
)

var (
	ErrImported     = Define(BlameThirdParty, NamespaceRuntime, ReasonInternal)
	ErrDecodingJSON = Define(BlameThirdParty, NamespaceRuntime, ReasonRequestDecoding)
	ErrFormat       = Define(BlameThirdParty, NamespaceRuntime, ReasonRequestFormat)
)

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

func (e *Error) UnmarshalJSON(data []byte) error {
	var errJson errorJSON

	err := json.Unmarshal(data, &errJson)
	if err != nil {
		return ErrDecodingJSON.Wrap(err)
	}

	code := strings.SplitN(errJson.Code, ".", 3)
	if len(code) != 3 {
		return ErrFormat.YeetExplain("code must be formed of 3 parts")
	}

	e.blame = mapCodeToBlame[code[0]]
	e.namespace = mapCodeToNamespace[code[1]]
	e.reason = mapCodeToReason[code[2]]

	e.explanation.WriteString(errJson.Explain)
	e.help = errJson.Help
	e.defined = ErrImported
	e.multi = errJson.Multi

	return nil
}
