package oops_json

import (
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

type testCustomTypeNothing struct {
}

type testCustomTypeDecodeFloatString float64

var testCustomTypeDecodeFloatStringErr = oops.Define(oops.BlameDeveloper, oops.NamespaceTest, oops.ReasonResourceDecoding)

func (t *testCustomTypeDecodeFloatString) UnmarshalJSON(b []byte) error {
	s := string(b)
	if s == `"oops_error"` {
		return testCustomTypeDecodeFloatStringErr.YeetExplain("failed on purpose")
	}

	if s[0] == '"' && s[len(s)-1] == '"' {
		s = s[1 : len(s)-1]
	}

	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}

	*t = testCustomTypeDecodeFloatString(f)
	return nil
}

func TestError_decode(t *testing.T) {
	t.Parallel()

	t.Run("nil error", func(t *testing.T) {
		var out interface{}
		input := `{"foo": "bar"}`
		err := Error(json.Unmarshal([]byte(input), &out))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("invalid json", func(t *testing.T) {
		var out interface{}
		input := `"foo": "bar"`
		err := Error(json.Unmarshal([]byte(input), &out))
		if err == nil {
			t.Fatal("expected error")
		}

		if !errors.Is(err, ErrInvalid) {
			t.Fatalf("expected error %v, got %v", ErrInvalid, err)
		}
		t.Logf("%+v", err)
	})
	t.Run("xml", func(t *testing.T) {
		var out interface{}
		input := `<data> foo bar baz </data>`
		err := Error(json.Unmarshal([]byte(input), &out))
		if err == nil {
			t.Fatal("expected error")
		}

		if !errors.Is(err, ErrInvalid) {
			t.Fatalf("expected error %v, got %v", ErrInvalid, err)
		}
		t.Logf("%+v", err)
	})
	t.Run("typed", func(t *testing.T) {
		t.Run("decode string as float", func(t *testing.T) {
			var out struct {
				Name  string  `json:"name"`
				Value float64 `json:"value"`
			}

			input := `{"name": "foo", "value": "3.1415"}`
			err := Error(json.Unmarshal([]byte(input), &out))
			if err == nil {
				t.Fatal("expected error")
			}

			if !errors.Is(err, ErrDecoding) {
				t.Fatalf("expected error %v, got %v", ErrDecoding, err)
			}
			t.Logf("%+v", err)

			wantStr := "CLIENT.RUNTIME.RESOURCE_DECODING(check byte at index=33, field='value', type expected='float64' got='string')"
			if err.Error() != wantStr {
				t.Fatalf("expected error %q, got %q", wantStr, err.Error())
			}
		})
		t.Run("decode int as float", func(t *testing.T) {
			var out struct {
				Name  string  `json:"name"`
				Value float64 `json:"value"`
			}
			input := `{"name": "foo", "value": 3}`
			err := Error(json.Unmarshal([]byte(input), &out))
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
		})
		t.Run("decode float as int", func(t *testing.T) {
			var out struct {
				Name  string `json:"name"`
				Value int    `json:"value"`
			}
			input := `{"name": "foo", "value": 3.1415}`
			err := Error(json.Unmarshal([]byte(input), &out))
			if err == nil {
				t.Fatal("expected error")
			}
			t.Logf("%+v", err)

			wantStr := "CLIENT.RUNTIME.RESOURCE_DECODING(check byte at index=31, field='value', type expected='int' got='number 3.1415')"
			if err.Error() != wantStr {
				t.Fatalf("expected error %q, got %q", wantStr, err.Error())
			}
		})
	})
}

func TestError_decode_customTypes(t *testing.T) {
	t.Parallel()

	t.Run("nothing and extra unused field", func(t *testing.T) {
		var out struct {
			Nothing testCustomTypeNothing `json:"nothing"`
		}

		input := `{"nothing": {}, "custom": {"value": 3.1415}}`
		err := Error(json.Unmarshal([]byte(input), &out))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
	})
	t.Run("unmarshal ok", func(t *testing.T) {
		var out struct {
			Custom testCustomTypeDecodeFloatString `json:"custom"`
		}

		input := `{"custom": "3.1415"}`
		err := Error(json.Unmarshal([]byte(input), &out))
		if err != nil {
			t.Errorf("unexpected error: %s", err)
		}
	})
	t.Run("unmarshal err", func(t *testing.T) {
		var out struct {
			Custom testCustomTypeDecodeFloatString `json:"custom"`
		}

		input := `{"custom": "foobar"}`
		err := Error(json.Unmarshal([]byte(input), &out))
		if err == nil {
			t.Fatal("expected error")
		}
		t.Logf("%+v", err)

		wantStr := "DEVELOPER.RUNTIME.UNEXPECTED(unexpected json error)"
		if err.Error() != wantStr {
			t.Fatalf("expected error %q, got %q", wantStr, err.Error())
		}
	})
	t.Run("custom oops error", func(t *testing.T) {
		var out struct {
			Custom testCustomTypeDecodeFloatString `json:"custom"`
		}

		input := `{"custom": "oops_error"}`
		err := Error(json.Unmarshal([]byte(input), &out))
		if err == nil {
			t.Fatal("expected error")
		}
		t.Logf("%+v", err)

		wantStr := "DEVELOPER.TEST.RESOURCE_DECODING(failed on purpose, json error)"
		if err.Error() != wantStr {
			t.Fatalf("expected error %q, got %q", wantStr, err.Error())
		}
	})
}

func TestError_encode(t *testing.T) {
	t.Parallel()

	t.Run("map string interface{}", func(t *testing.T) {
		input := map[string]interface{}{
			"name":  "bar",
			"value": 3.141516,
		}

		out, err := ErrorM(json.Marshal(input))
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		wantStr := `{"name":"bar","value":3.141516}`
		if string(out) != wantStr {
			t.Fatalf("expected %q, got %q", wantStr, string(out))
		}
	})
	t.Run("bad type", func(t *testing.T) {
		var input struct {
			Name    string `json:"name"`
			FooFunc func() `json:"foo_func"`
		}

		_, err := ErrorM(json.Marshal(input))
		if err == nil {
			t.Fatal("expected error")
		}
		t.Logf("%+v", err)

		wantStr := "CLIENT.RUNTIME.RESOURCE_ENCODING(unsupported type='func()')"
		if err.Error() != wantStr {
			t.Fatalf("expected error %q, got %q", wantStr, err.Error())
		}
	})
	t.Run("bad value", func(t *testing.T) {
		var input struct {
			Name  string  `json:"name"`
			Value float64 `json:"value"`
		}

		input.Value = math.NaN()

		_, err := ErrorM(json.Marshal(input))
		if err == nil {
			t.Fatal("expected error")
		}
		t.Logf("%+v", err)

		wantStr := "CLIENT.RUNTIME.RESOURCE_ENCODING(unsupported value type='float64' string='NaN')"
		if err.Error() != wantStr {
			t.Fatalf("expected error %q, got %q", wantStr, err.Error())
		}
	})
}

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("special", func(t *testing.T) {
		t.Run("unwrap MarshalerError nil", func(t *testing.T) {
			err := Error(new(json.MarshalerError))
			if err == nil {
				t.Fatal("expected error")
			}
			t.Logf("%+v", err)

			wantStr := "DEVELOPER.RUNTIME.UNEXPECTED(unexpected nested marshaler error)"
			if err.Error() != wantStr {
				t.Fatalf("expected error %q, got %q", wantStr, err.Error())
			}
		})
		t.Run("unwrap MarshalerError ok", func(t *testing.T) {
			mErr := &json.MarshalerError{
				Err: errors.New("foo bar"),
			}

			err := Error(mErr)
			if err == nil {
				t.Fatal("expected error")
			}
			t.Logf("%+v", err)

			wantStr := "DEVELOPER.RUNTIME.UNEXPECTED(unexpected json error)"
			if err.Error() != wantStr {
				t.Fatalf("expected error %q, got %q", wantStr, err.Error())
			}
		})
	})
}
