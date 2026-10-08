package oops_test

import (
	"errors"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

func TestError_Explain(t *testing.T) {
	t.Parallel()

	t.Run("single explanation", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Explainf("something happened")
		if err.Explanation() != "something happened" {
			t.Fatalf("got %q", err.Explanation())
		}
	})

	t.Run("multiple explanations", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Explainf("first").Explainf("second")
		if err.Explanation() != "first, second" {
			t.Fatalf("got %q", err.Explanation())
		}
	})

	t.Run("Explainf", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Explainf("count=%d", 5)
		if err.Explanation() != "count=5" {
			t.Fatalf("got %q", err.Explanation())
		}
	})

	t.Run("escaped percent without args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeetf("100%%").Explainf("50%% done")
		if err.Explanation() != "100%, 50% done" {
			t.Fatalf("got %q", err.Explanation())
		}
	})

	t.Run("empty explanation is skipped", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Explainf("first").Explainf("").Explainf("third")
		if err.Explanation() != "first, third" {
			t.Fatalf("got %q", err.Explanation())
		}
	})
}

func TestError_CausesActions(t *testing.T) {
	t.Parallel()

	t.Run("AddCause", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Causes(oops.CauseNotFound).Yeet()
		err = err.AddCauses(oops.CauseTimeout)
		if !err.HasCause(oops.CauseNotFound) {
			t.Fatal("missing CauseNotFound")
		}
		if !err.HasCause(oops.CauseTimeout) {
			t.Fatal("missing CauseTimeout")
		}
	})

	t.Run("SetActions keeps its own copy of the actions", func(t *testing.T) {
		t.Parallel()
		actions := []oops.Action{oops.ActionRetry}
		err := oops.Define("test").Yeet().SetActions(actions...)
		actions[0] = oops.ActionAbort
		if got := err.Actions(); len(got) != 1 || got[0] != oops.ActionRetry {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("SetActions replaces", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Actions(oops.ActionRetry).Yeet()
		_ = err.SetActions(oops.ActionAbort)
		if err.HasAction(oops.ActionRetry) {
			t.Fatal("ActionRetry should have been replaced")
		}
		if !err.HasAction(oops.ActionAbort) {
			t.Fatal("missing ActionAbort")
		}
	})
}

func TestError_Fields(t *testing.T) {
	t.Parallel()

	t.Run("Set and Get", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Set("key", "value")
		v, ok := err.Get("key")
		if !ok || v != "value" {
			t.Fatalf("got %v, %v", v, ok)
		}
	})

	t.Run("Get missing key", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		_, ok := err.Get("missing")
		if ok {
			t.Fatal("expected false for missing key")
		}
	})

	t.Run("Fields len", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Set("a", 1).Set("b", 2)
		all := err.Fields()
		if len(all) != 2 {
			t.Fatalf("expected 2 fields, got %d", len(all))
		}
	})
}

func TestError_Nest(t *testing.T) {
	t.Parallel()

	t.Run("Nest adds to wrapped", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("inner")
		err := oops.Define("test").Yeet().Nest(inner)
		if len(err.Unwrap()) != 1 {
			t.Fatalf("expected 1 wrapped, got %d", len(err.Unwrap()))
		}
	})

	t.Run("Nest several errors", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		parent := def.Yeet()
		child1 := def.Yeet()
		child2 := def.Yeet()
		_ = parent.Nest(child1, child2)
		if len(parent.Unwrap()) != 2 {
			t.Fatalf("expected 2 wrapped, got %d", len(parent.Unwrap()))
		}
	})

	t.Run("Nest skips nil", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		err := oops.Define("test").Yeet()
		_ = err.Nest(nil, typedNil)
		if len(err.Unwrap()) != 0 {
			t.Fatal("expected 0 wrapped")
		}
	})
}
