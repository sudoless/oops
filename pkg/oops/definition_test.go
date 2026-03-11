package oops_test

import (
	"errors"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestDefine(t *testing.T) {
	t.Parallel()

	t.Run("creates definition with code", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test.error")
		if def.Code() != "test.error" {
			t.Fatalf("expected code %q, got %q", "test.error", def.Code())
		}
	})

	t.Run("Error returns code", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test.error")
		if def.Error() != "test.error" {
			t.Fatalf("expected %q, got %q", "test.error", def.Error())
		}
	})

	t.Run("Error returns code and message", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test.error").Message("something went wrong")
		if def.Error() != "test.error: something went wrong" {
			t.Fatalf("expected %q, got %q", "test.error: something went wrong", def.Error())
		}
	})
}

func TestErrorDefinition_Is(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		if def.Is(nil) {
			t.Fatal("expected false for nil")
		}
	})

	t.Run("same definition", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		if !def.Is(def) {
			t.Fatal("expected true for same definition")
		}
	})

	t.Run("different definition", func(t *testing.T) {
		t.Parallel()
		def1 := oops.Define("test1")
		def2 := oops.Define("test2")
		if def1.Is(def2) {
			t.Fatal("expected false for different definitions")
		}
	})

	t.Run("error from this definition", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if !def.Is(err) {
			t.Fatal("expected true for error from this definition")
		}
	})

	t.Run("error from different definition", func(t *testing.T) {
		t.Parallel()
		def1 := oops.Define("test1")
		def2 := oops.Define("test2")
		err := def2.Yeet()
		if def1.Is(err) {
			t.Fatal("expected false for error from different definition")
		}
	})

	t.Run("inherits", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base")
		child := oops.Define("child").Inherits(base)
		if !child.Is(base) {
			t.Fatal("expected child.Is(base) to be true")
		}
		if base.Is(child) {
			t.Fatal("expected base.Is(child) to be false")
		}
	})

	t.Run("errors.Is with definition target", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if !errors.Is(err, def) {
			t.Fatal("expected errors.Is(err, def) to be true")
		}
	})

	t.Run("errors.Is with inherited definition", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base")
		child := oops.Define("child").Inherits(base)
		err := child.Yeet()
		if !errors.Is(err, base) {
			t.Fatal("expected errors.Is(err, base) to be true via inheritance")
		}
	})
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if err == nil {
			t.Fatal("expected non-nil error")
		}
		if err.Code() != "test" {
			t.Fatalf("expected code %q, got %q", "test", err.Code())
		}
	})

	t.Run("with explanation", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeetf("something %s", "bad")
		if err.Explanation() != "something bad" {
			t.Fatalf("expected explanation %q, got %q", "something bad", err.Explanation())
		}
	})

	t.Run("copies causes", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test").Causes(oops.CauseNotFound)
		err := def.Yeet()
		if !err.HasCause(oops.CauseNotFound) {
			t.Fatal("expected error to have CauseNotFound")
		}
	})

	t.Run("copies actions", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test").Actions(oops.ActionRetry)
		err := def.Yeet()
		if !err.HasAction(oops.ActionRetry) {
			t.Fatal("expected error to have ActionRetry")
		}
	})
}

func TestYeet(t *testing.T) {
	t.Parallel()

	t.Run("same as New", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if err.Code() != "test" {
			t.Fatalf("expected code %q, got %q", "test", err.Code())
		}
	})

	t.Run("Yeetf with format", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeetf("failed: %d", 42)
		if err.Explanation() != "failed: 42" {
			t.Fatalf("expected explanation %q, got %q", "failed: 42", err.Explanation())
		}
	})
}

func TestWrap(t *testing.T) {
	t.Parallel()

	t.Run("wraps error", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		inner := errors.New("inner")
		err := def.Wrap(inner)
		wrapped := err.Unwrap()
		if len(wrapped) != 1 || wrapped[0] != inner {
			t.Fatal("expected inner error to be wrapped")
		}
	})

	t.Run("wrap nil", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Wrap(nil)
		if len(err.Unwrap()) != 0 {
			t.Fatal("expected no wrapped errors for nil")
		}
	})

	t.Run("Wrapf with format", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		inner := errors.New("inner")
		err := def.Wrapf(inner, "wrapping: %s", "context")
		if err.Explanation() != "wrapping: context" {
			t.Fatalf("expected explanation %q, got %q", "wrapping: context", err.Explanation())
		}
		if len(err.Unwrap()) != 1 {
			t.Fatal("expected one wrapped error")
		}
	})

	t.Run("errors.Is finds wrapped error", func(t *testing.T) {
		t.Parallel()
		inner := oops.Define("inner")
		outer := oops.Define("outer")
		innerErr := inner.Yeet()
		outerErr := outer.Wrap(innerErr)
		if !errors.Is(outerErr, inner) {
			t.Fatal("expected errors.Is to find inner through wrapping")
		}
	})
}

func TestTrace(t *testing.T) {
	t.Parallel()

	t.Run("no trace by default", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if len(err.Trace()) != 0 {
			t.Fatal("expected no trace")
		}
	})

	t.Run("trace when Traced", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test").Traced()
		err := def.Yeet()
		if len(err.Trace()) == 0 {
			t.Fatal("expected trace to be captured")
		}
	})
}
