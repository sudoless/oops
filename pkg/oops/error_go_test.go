package oops_test

import (
	"errors"
	"fmt"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

// customError is a concrete error type used to test As/Is with wrapped non-oops errors.
type customError struct{ code int }

func (e *customError) Error() string { return fmt.Sprintf("custom(%d)", e.code) }

func TestError_ErrorString(t *testing.T) {
	t.Parallel()

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()
		var err *oops.Error
		if err.Error() != "oops.Error(nil)" {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("zero value", func(t *testing.T) {
		t.Parallel()
		err := &oops.Error{}
		if err.Error() != "oops.Error(undefined)" {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("code only", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		if err.Error() != "test" {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("code with message", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Message("msg").Yeet()
		if err.Error() != "test: msg" {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("code with explanation", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeetf("explained")
		if err.Error() != "test: explained" {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("code with explanation and message", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Message("msg").Yeetf("explained")
		if err.Error() != "test: msg; explained" {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("custom formatter", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test").Formatter(func(err *oops.Error) string {
			return "custom: " + err.Explanation()
		})
		err := def.Yeetf("hello")
		if err.Error() != "custom: hello" {
			t.Fatalf("got %q", err.Error())
		}
	})
}

func TestError_Unwrap(t *testing.T) {
	t.Parallel()

	t.Run("returns wrapped slice", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("inner")
		err := oops.Define("test").Wrap(inner)
		unwrapped := err.Unwrap()
		if len(unwrapped) != 1 || !errors.Is(unwrapped[0], inner) {
			t.Fatal("expected inner error")
		}
	})

	t.Run("empty when no wrapping", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		if len(err.Unwrap()) != 0 {
			t.Fatal("expected empty")
		}
	})
}

func TestError_Is(t *testing.T) {
	t.Parallel()

	t.Run("nil error is nil", func(t *testing.T) {
		t.Parallel()
		var err *oops.Error
		var typedNil *oops.Error
		if !err.Is(nil) || !err.Is(typedNil) {
			t.Fatal("nil error should Is nil")
		}
	})

	t.Run("typed-nil target", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		err := oops.Define("test").Yeet()
		if err.Is(typedNil) {
			t.Fatal("expected false for typed-nil target")
		}
		if errors.Is(err, typedNil) {
			t.Fatal("expected errors.Is false for typed-nil target")
		}
	})

	t.Run("matches definition", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if !err.Is(def) {
			t.Fatal("expected match")
		}
	})

	t.Run("matches same-def error", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err1 := def.Yeet()
		err2 := def.Yeet()
		if !err1.Is(err2) {
			t.Fatal("expected match for same definition")
		}
	})

	t.Run("no match different def", func(t *testing.T) {
		t.Parallel()
		def1 := oops.Define("a")
		def2 := oops.Define("b")
		if def1.Yeet().Is(def2) {
			t.Fatal("expected no match")
		}
		if errors.Is(def1.Yeet(), def2) {
			t.Fatal("expected errors.Is no match")
		}
	})

	t.Run("compares definitions only, not wrapped errors", func(t *testing.T) {
		t.Parallel()
		inner := oops.Define("inner")
		err := oops.Define("outer").Wrap(inner.Yeet())
		if err.Is(inner) {
			t.Fatal("(*Error).Is must not traverse")
		}
	})

	t.Run("errors.Is traverses wrapped", func(t *testing.T) {
		t.Parallel()
		inner := oops.Define("inner")
		outer := oops.Define("outer")
		err := outer.Wrap(inner.Yeet())
		if !errors.Is(err, inner) {
			t.Fatal("errors.Is should find inner via unwrap")
		}
	})

	t.Run("errors.Is with stdlib error wrapped", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("sentinel")
		err := oops.Define("test").Wrap(sentinel)
		if !errors.Is(err, sentinel) {
			t.Fatal("should find sentinel via unwrap")
		}
	})

	t.Run("errors.Join without the target", func(t *testing.T) {
		t.Parallel()
		target := oops.Define("target")
		joined := errors.Join(errors.New("other"), oops.Define("unrelated").Yeet())
		if errors.Is(joined, target) {
			t.Fatal("expected no match")
		}
	})
}

func TestError_Is_InheritanceAsymmetry(t *testing.T) {
	t.Parallel()

	base := oops.Define("base")
	child := oops.Define("child").Inherits(base)

	childErr := child.Yeet()
	baseErr := base.Yeet()

	if !childErr.Is(baseErr) {
		t.Error("childErr.Is(baseErr) should be true: same inheritance contract as Is(def)")
	}

	if baseErr.Is(childErr) {
		t.Error("baseErr.Is(childErr) should be false: base does not inherit child")
	}

	if errors.Is(baseErr, child) {
		t.Error("errors.Is(baseErr, child) should be false")
	}
}

func TestError_As(t *testing.T) {
	t.Parallel()

	t.Run("errors.As extracts Error", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeetf("hello")
		var target *oops.Error
		if !errors.As(err, &target) {
			t.Fatal("expected As to succeed")
		}
		if target.Explanation() != "hello" {
			t.Fatalf("got %q", target.Explanation())
		}
	})

	t.Run("errors.As through wrapping returns outer first", func(t *testing.T) {
		t.Parallel()
		inner := oops.Define("inner").Yeetf("deep")
		outer := oops.Define("outer").Wrap(inner)

		var target *oops.Error
		if !errors.As(outer, &target) {
			t.Fatal("expected As to succeed")
		}
		if target.Code() != "outer" {
			t.Fatalf("got %q", target.Code())
		}
	})

	t.Run("errors.As extracts wrapped non-oops error", func(t *testing.T) {
		t.Parallel()
		wrapped := &customError{code: 42}
		err := oops.Define("test").Wrap(wrapped)
		var target *customError
		if !errors.As(err, &target) {
			t.Fatal("errors.As should find wrapped custom error via Unwrap")
		}
		if target.code != 42 {
			t.Fatalf("got code %d", target.code)
		}
	})

	t.Run("errors.As extracts deeply wrapped non-oops error", func(t *testing.T) {
		t.Parallel()
		sentinel := &customError{code: 99}
		middle := fmt.Errorf("wrap: %w", sentinel)
		err := oops.Define("test").Wrap(middle)
		var target *customError
		if !errors.As(err, &target) {
			t.Fatal("errors.As should reach deeply wrapped custom error")
		}
		if target.code != 99 {
			t.Fatalf("got code %d", target.code)
		}
	})
}
