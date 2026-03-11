package oops_test

import (
	"errors"
	"fmt"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

// customErr is a concrete error type used to test As/Is with wrapped non-oops errors.
type customErr struct{ code int }

func (e *customErr) Error() string { return fmt.Sprintf("custom(%d)", e.code) }

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
		_ = err.AddCause(oops.CauseTimeout)
		if !err.HasCause(oops.CauseNotFound) {
			t.Fatal("missing CauseNotFound")
		}
		if !err.HasCause(oops.CauseTimeout) {
			t.Fatal("missing CauseTimeout")
		}
	})

	t.Run("WithActions replaces", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Actions(oops.ActionRetry).Yeet()
		_ = err.WithActions(oops.ActionAbort)
		if err.HasAction(oops.ActionRetry) {
			t.Fatal("ActionRetry should have been replaced")
		}
		if !err.HasAction(oops.ActionAbort) {
			t.Fatal("missing ActionAbort")
		}
	})
}

func TestError_Path(t *testing.T) {
	t.Parallel()

	t.Run("WithPathf with args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().WithPathf("user/%d", 42)
		if err.Path() != "user/42" {
			t.Fatalf("got %q", err.Path())
		}

		args := err.PathArgs()
		if len(args) != 1 || len(args) != 1 {
			t.Fatalf("expected 1 path arg, got %v", args)
		}
	})

	t.Run("WithPath sets nil args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().WithPathf("static")
		args := err.PathArgs()
		if args != nil {
			t.Fatalf("expected nil args, got %v", args)
		}
	})

	t.Run("WithPathf no args sets nil args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		err.WithPathf("static") // return value not captured; mutation is on receiver
		args := err.PathArgs()
		if args != nil {
			t.Fatalf("expected nil args, got %v", args)
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

func TestError_Accessors(t *testing.T) {
	t.Parallel()

	t.Run("Definition", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if err.Definition() != def {
			t.Fatal("expected definition to match")
		}
	})

	t.Run("Code", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("my.code").Yeet()
		if err.Code() != "my.code" {
			t.Fatalf("got %q", err.Code())
		}
	})

	t.Run("Message", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Message("hello").Yeet()
		if err.Message() != "hello" {
			t.Fatalf("got %q", err.Message())
		}
	})
}

func TestError_ErrorString(t *testing.T) {
	t.Parallel()

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()
		var err *oops.Error
		if err.Error() != "oops.Error(nil)" {
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
		def := oops.Define("test").SetFormatter(func(err *oops.Error) string {
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
		if len(unwrapped) != 1 || unwrapped[0] != inner {
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
		if !err.Is(nil) {
			t.Fatal("nil error should Is nil")
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
	})

	t.Run("inherits chain", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base")
		child := oops.Define("child").Inherits(base)
		err := child.Yeet()
		if !err.Is(base) {
			t.Fatal("expected match via inheritance")
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

	t.Run("direct Is with stdlib error wrapped", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("sentinel")
		err := oops.Define("test").Wrap(sentinel)
		if !err.Is(sentinel) {
			t.Fatal("Is should match wrapped stdlib error")
		}
	})

	t.Run("direct Is with deeply wrapped stdlib error", func(t *testing.T) {
		t.Parallel()
		sentinel := errors.New("sentinel")
		middle := fmt.Errorf("middle: %w", sentinel)
		err := oops.Define("test").Wrap(middle)
		if !err.Is(sentinel) {
			t.Fatal("Is should match deeply wrapped stdlib error")
		}
	})
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
		wrapped := &customErr{code: 42}
		err := oops.Define("test").Wrap(wrapped)
		var target *customErr
		if !errors.As(err, &target) {
			t.Fatal("errors.As should find wrapped custom error via Unwrap")
		}
		if target.code != 42 {
			t.Fatalf("got code %d", target.code)
		}
	})

	t.Run("errors.As extracts deeply wrapped non-oops error", func(t *testing.T) {
		t.Parallel()
		sentinel := &customErr{code: 99}
		middle := fmt.Errorf("wrap: %w", sentinel)
		err := oops.Define("test").Wrap(middle)
		var target *customErr
		if !errors.As(err, &target) {
			t.Fatal("errors.As should reach deeply wrapped custom error")
		}
		if target.code != 99 {
			t.Fatalf("got code %d", target.code)
		}
	})
}
