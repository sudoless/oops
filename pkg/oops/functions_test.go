package oops_test

import (
	"errors"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestCatch(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		if oops.Catch(nil) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("oops error passes through", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		caught := oops.Catch(err)
		if caught != err {
			t.Fatal("expected same error")
		}
	})

	t.Run("stdlib error wrapped with ErrUncaught", func(t *testing.T) {
		t.Parallel()
		err := errors.New("plain")
		caught := oops.Catch(err)
		if !errors.Is(caught, oops.ErrUncaught) {
			t.Fatal("expected ErrUncaught wrapping")
		}
		if !errors.Is(caught, err) {
			t.Fatal("should still unwrap to original")
		}
	})
}

func TestExplainf(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		if oops.Explainf(nil, "ignored %d", 42) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("formats explanation", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		result := oops.Explainf(err, "count=%d", 5)
		if result.Explanation() != "count=5" {
			t.Fatalf("got %q", result.Explanation())
		}
	})
}

func TestAddCause(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		if oops.AddCause(nil, oops.CauseAuth) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("adds cause", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		result := oops.AddCause(err, oops.CauseAuth)
		if !result.HasCause(oops.CauseAuth) {
			t.Fatal("expected CauseAuth")
		}
	})
}

func TestPathf(t *testing.T) {
	t.Parallel()

	t.Run("formats path segment", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		result := oops.Pathf(err, "user/%d", 42)
		if result.Path() != "user/42" {
			t.Fatalf("got %q", result.Path())
		}
	})
}

func TestAs(t *testing.T) {
	t.Parallel()

	t.Run("nil returns false", func(t *testing.T) {
		t.Parallel()
		_, ok := oops.As(nil, oops.Define("test"))
		if ok {
			t.Fatal("expected false")
		}
	})

	t.Run("nil target returns false", func(t *testing.T) {
		t.Parallel()
		_, ok := oops.As(oops.Define("test").Yeet(), nil)
		if ok {
			t.Fatal("expected false")
		}
	})

	t.Run("direct match", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		found, ok := oops.As(err, def)
		if !ok || found != err {
			t.Fatal("expected direct match")
		}
	})

	t.Run("wrapped match", func(t *testing.T) {
		t.Parallel()
		inner := oops.Define("inner")
		outer := oops.Define("outer")
		innerErr := inner.Yeetf("deep")
		outerErr := outer.Wrap(innerErr)

		found, ok := oops.As(outerErr, inner)
		if !ok {
			t.Fatal("expected match")
		}
		if found.Explanation() != "deep" {
			t.Fatalf("got %q", found.Explanation())
		}
	})

	t.Run("via stdlib wrapping", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("target")
		innerErr := def.Yeetf("inner")
		joined := errors.Join(errors.New("other"), innerErr)

		found, ok := oops.As(joined, def)
		if !ok {
			t.Fatal("expected match through errors.Join")
		}
		if found.Explanation() != "inner" {
			t.Fatalf("got %q", found.Explanation())
		}
	})

	t.Run("inherits match", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base")
		child := oops.Define("child").Inherits(base)
		err := child.Yeet()

		found, ok := oops.As(err, base)
		if !ok {
			t.Fatal("expected match via inheritance")
		}
		if found.Code() != "child" {
			t.Fatalf("expected child code, got %q", found.Code())
		}
	})
}

func TestNest(t *testing.T) {
	t.Parallel()

	t.Run("nil def returns nil", func(t *testing.T) {
		t.Parallel()
		if oops.Nest(nil, errors.New("err")) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("no errors returns nil", func(t *testing.T) {
		t.Parallel()
		if oops.Nest(oops.Define("test")) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("all nil errors returns nil", func(t *testing.T) {
		t.Parallel()
		if oops.Nest(oops.Define("test"), nil, nil) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("creates parent with wrapped children", func(t *testing.T) {
		t.Parallel()
		parent := oops.Define("parent")
		err1 := oops.Define("child1").Yeet()
		err2 := oops.Define("child2").Yeet()

		result := oops.Nest(parent, err1, err2)
		if result == nil {
			t.Fatal("expected non-nil")
		}
		if result.Code() != "parent" {
			t.Fatalf("expected parent code, got %q", result.Code())
		}
		if len(result.Unwrap()) != 2 {
			t.Fatalf("expected 2 wrapped, got %d", len(result.Unwrap()))
		}
	})

	t.Run("skips nil errors in list", func(t *testing.T) {
		t.Parallel()
		parent := oops.Define("parent")
		err1 := oops.Define("child").Yeet()

		result := oops.Nest(parent, nil, err1, nil)
		if len(result.Unwrap()) != 1 {
			t.Fatalf("expected 1 wrapped, got %d", len(result.Unwrap()))
		}
	})
}

func TestPresets(t *testing.T) {
	t.Parallel()

	t.Run("ErrUncaught", func(t *testing.T) {
		t.Parallel()
		err := oops.ErrUncaught.Yeet()
		if err.Code() != "uncaught" {
			t.Fatalf("got %q", err.Code())
		}
		if !err.HasCause(oops.CauseInternal) {
			t.Fatal("expected CauseInternal")
		}
		if !err.HasAction(oops.ActionAbort) {
			t.Fatal("expected ActionAbort")
		}
		if len(err.Trace()) == 0 {
			t.Fatal("expected trace")
		}
	})

	t.Run("ErrTODO", func(t *testing.T) {
		t.Parallel()
		err := oops.ErrTODO.Yeet()
		if err.Code() != "todo" {
			t.Fatalf("got %q", err.Code())
		}
		if err.Message() != "not implemented" {
			t.Fatalf("got %q", err.Message())
		}
		if len(err.Trace()) == 0 {
			t.Fatal("expected trace")
		}
	})
}
