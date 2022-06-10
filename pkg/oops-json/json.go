package oops_json

import (
	"encoding/json"
	"go.sdls.io/oops/pkg/oops"
)

// Error can be used on any JSON (standard package) error to be wrapped inside oops.Error. This is useful as the
// wrapping moves certain error details in the oops.Error explanation and also hides internal struct details and
// keeps the wrapped error for later access.
func Error(jsonErr error) error {
	if jsonErr == nil {
		return nil
	}

	switch v := jsonErr.(type) {
	case *json.MarshalerError:
		un := v.Unwrap()
		if un == nil {
			return oops.ErrUnexpected.Wrap(jsonErr, "unexpected nested marshaler error")
		}
		jsonErr = un
	}

	switch v := jsonErr.(type) {
	case *json.SyntaxError:
		return ErrInvalid.Wrap(jsonErr, "check byte at index=%d", v.Offset)
	case *json.UnmarshalTypeError:
		return ErrDecoding.Wrap(jsonErr, "check byte at index=%d field='%s' type expected='%s' got='%s'",
			v.Offset, v.Field, v.Type.String(), v.Value)
	case *json.UnsupportedTypeError:
		return ErrEncoding.Wrap(jsonErr, "unsupported type='%s'", v.Type.String())
	case *json.UnsupportedValueError:
		return ErrEncoding.Wrap(jsonErr, "unsupported value type='%s' string='%s'",
			v.Value.Type().String(), v.Str)
	default:
		if jsonErr.Error() == "unexpected EOF" {
			return ErrInvalid.Wrap(jsonErr, "unexpected end of JSON")
		}

		oopsErr, is, _ := oops.As(jsonErr)
		if is {
			return oops.Explain(oopsErr, "json error")
		}

		return oops.ErrUnexpected.Wrap(jsonErr, "unexpected json error")
	}
}

// ErrorM is a helper function which can be used around json.Marshal to wrap the error inside oops.Error and return
// the unmodified output from json.Marshal.
// Example
// out, err := oopsJSON.ErrorM(json.Marshal(&input))
func ErrorM(out []byte, err error) ([]byte, error) {
	return out, Error(err)
}
