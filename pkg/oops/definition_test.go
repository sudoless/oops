package oops_test

import (
	"errors"
	"slices"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
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

	t.Run("Error panics", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test.error").Message("something went wrong")
		defer func() {
			if recover() == nil {
				t.Fatal("expected definition Error() to panic")
			}
		}()
		_ = def.Error()
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

	t.Run("typed-nil *Error", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		if oops.Define("test").Is(typedNil) {
			t.Fatal("expected false for typed-nil *Error")
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
		if def.Is(err) {
			t.Fatal("expected false: a definition only matches definitions")
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

func TestDefinition_Builders(t *testing.T) {
	t.Parallel()

	t.Run("receiver unchanged", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base").Causes(oops.CauseIO).Actions(oops.ActionRetry)
		parent := oops.Define("parent")

		derived := []*oops.ErrorDefinition{
			base.Causes(oops.CauseTimeout),
			base.Actions(oops.ActionAbort),
			base.Message("changed"),
			base.Traced(),
			base.Inherits(parent),
			base.Formatter(func(*oops.Error) string { return "formatted" }),
		}
		for idx, built := range derived {
			if built == base {
				t.Fatalf("builder %d returned the receiver", idx)
			}
		}

		err := base.Yeet()
		if got := err.Causes(); !slices.Equal(got, []string{oops.CauseIO}) {
			t.Errorf("causes = %v", got)
		}
		if got := err.Actions(); !slices.Equal(got, []string{oops.ActionRetry}) {
			t.Errorf("actions = %v", got)
		}
		if err.Message() != "" {
			t.Errorf("message = %q", err.Message())
		}
		if err.Trace() != nil {
			t.Error("receiver became traced")
		}
		if errors.Is(err, parent) {
			t.Error("receiver gained a parent")
		}
		if err.Error() != "base" {
			t.Errorf("Error() = %q", err.Error())
		}
	})

	t.Run("derived definition has receiver's tags plus its own", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base").Causes(oops.CauseIO)
		derived := base.Causes(oops.CauseTimeout)
		if got := derived.Yeet().Causes(); !slices.Equal(got, []string{oops.CauseIO, oops.CauseTimeout}) {
			t.Fatalf("causes = %v", got)
		}
	})

	t.Run("sibling builders do not share tags", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base").Causes(oops.CauseIO, oops.CauseTimeout)
		left := base.Causes(oops.CauseAuth)
		right := base.Causes(oops.CauseExpired)
		if got := left.Yeet().Causes(); !slices.Equal(got, []string{oops.CauseIO, oops.CauseTimeout, oops.CauseAuth}) {
			t.Fatalf("left causes = %v", got)
		}
		if got := right.Yeet().Causes(); !slices.Equal(got, []string{oops.CauseIO, oops.CauseTimeout, oops.CauseExpired}) {
			t.Fatalf("right causes = %v", got)
		}
	})

	t.Run("inherits cannot form a cycle", func(t *testing.T) {
		t.Parallel()
		a := oops.Define("a")
		b := oops.Define("b").Inherits(a)
		a2 := a.Inherits(b)
		if errors.Is(a.Yeet(), b) {
			t.Fatal("a must not inherit b: the builder returned a new definition")
		}
		if !errors.Is(a2.Yeet(), a) || !errors.Is(a2.Yeet(), b) {
			t.Fatal("a2 should match a (via b) and b")
		}
	})

	t.Run("preset unchanged", func(t *testing.T) {
		t.Parallel()
		_ = oops.ErrForeign.Message("changed").Causes(oops.CauseIO).Actions(oops.ActionRetry)
		err := oops.ErrForeign.Yeet()
		if err.Message() != "" {
			t.Errorf("message = %q", err.Message())
		}
		if got := err.Causes(); !slices.Equal(got, []string{oops.CauseInternal}) {
			t.Errorf("causes = %v", got)
		}
		if got := err.Actions(); !slices.Equal(got, []string{oops.ActionAbort}) {
			t.Errorf("actions = %v", got)
		}
	})
}

func TestYeet(t *testing.T) {
	t.Parallel()

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

	t.Run("tags are independent of the definition", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test").Causes(oops.CauseNotFound).Actions(oops.ActionRetry)
		first := def.Yeet()
		first.Causes()[0] = "mutated"
		first.Actions()[0] = "mutated"

		second := def.Yeet()
		if got := second.Causes(); !slices.Equal(got, []string{oops.CauseNotFound}) {
			t.Fatalf("causes = %v", got)
		}
		if got := second.Actions(); !slices.Equal(got, []string{oops.ActionRetry}) {
			t.Fatalf("actions = %v", got)
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
		if len(wrapped) != 1 || !errors.Is(wrapped[0], inner) {
			t.Fatal("expected inner error to be wrapped")
		}
	})

	t.Run("wrap nil", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		def := oops.Define("test")
		for name, inner := range map[string]error{"nil": nil, "typed nil": typedNil} {
			if n := len(def.Wrap(inner).Unwrap()); n != 0 {
				t.Errorf("Wrap(%s): %d wrapped", name, n)
			}
			if n := len(def.Wrapf(inner, "x").Unwrap()); n != 0 {
				t.Errorf("Wrapf(%s): %d wrapped", name, n)
			}
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

	t.Run("wrapped text stays out of Error", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Wrapf(errors.New("inner"), "ctx")
		if err.Error() != "test: ctx" {
			t.Fatalf("got %q", err.Error())
		}
	})
}
