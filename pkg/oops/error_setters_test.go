package oops_test

import (
	"errors"
	"testing"

	"go.sdls.io/oops/pkg/oops"
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
		err = err.AddCause(oops.CauseTimeout)
		if !err.HasCause(oops.CauseNotFound) {
			t.Fatal("missing CauseNotFound")
		}
		if !err.HasCause(oops.CauseTimeout) {
			t.Fatal("missing CauseTimeout")
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

	t.Run("Fields returns live map", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Set("x", 42)
		if err.Fields()["x"] != 42 {
			t.Fatal("Fields should contain set value")
		}
	})
}

func TestError_NestAppend(t *testing.T) {
	t.Parallel()

	t.Run("Nest adds to wrapped", func(t *testing.T) {
		t.Parallel()
		inner := errors.New("inner")
		err := oops.Define("test").Yeet().Nest(inner)
		if len(err.Unwrap()) != 1 {
			t.Fatalf("expected 1 wrapped, got %d", len(err.Unwrap()))
		}
	})

	t.Run("Nest nil is no-op", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Nest(nil)
		if len(err.Unwrap()) != 0 {
			t.Fatal("expected 0 wrapped")
		}
	})

	t.Run("Append typed errors", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		parent := def.Yeet()
		child1 := def.Yeet()
		child2 := def.Yeet()
		_ = parent.Append(child1, child2)
		if len(parent.Unwrap()) != 2 {
			t.Fatalf("expected 2 wrapped, got %d", len(parent.Unwrap()))
		}
	})

	t.Run("Append skips nil", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		_ = err.Append(nil, nil)
		if len(err.Unwrap()) != 0 {
			t.Fatal("expected 0 wrapped")
		}
	})
}
