package oops_test

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

func TestForeign(t *testing.T) {
	t.Parallel()

	t.Run("nil returns nil", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		if oops.Foreign(nil) != nil || oops.Foreign(typedNil) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("oops error passes through", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		caught := oops.Foreign(err)
		if caught != err {
			t.Fatal("expected same error")
		}
	})

	t.Run("stdlib error wrapped with ErrForeign", func(t *testing.T) {
		t.Parallel()
		err := errors.New("plain")
		caught := oops.Foreign(err)
		if caught.Definition() != oops.ErrForeign {
			t.Fatalf("expected ErrForeign wrapping, got %q", caught.Code())
		}
		if !errors.Is(caught, err) {
			t.Fatal("should still unwrap to original")
		}
	})

	t.Run("does not look inside a foreign wrapper", func(t *testing.T) {
		t.Parallel()
		inner := oops.Define("inner").Yeet()
		outer := fmt.Errorf("context: %w", inner)
		caught := oops.Foreign(outer)
		if caught == inner || caught.Definition() != oops.ErrForeign {
			t.Fatalf("expected a new ErrForeign wrapper, got %q", caught.Code())
		}
	})
}

func TestNative(t *testing.T) {
	t.Parallel()

	t.Run("oops error", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		got, ok := oops.Native(err)
		if !ok || got != err {
			t.Fatalf("got %v, %v", got, ok)
		}
	})

	t.Run("not an oops error", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		cases := map[string]error{
			"nil":             nil,
			"typed nil":       typedNil,
			"stdlib":          errors.New("plain"),
			"foreign wrapper": fmt.Errorf("context: %w", oops.Define("inner").Yeet()),
		}
		for name, err := range cases {
			if got, ok := oops.Native(err); got != nil || ok {
				t.Errorf("%s: got %v, %v", name, got, ok)
			}
		}
	})
}

func TestHelpers(t *testing.T) {
	t.Parallel()

	var typedNil *oops.Error
	helpers := map[string]func(error) error{
		"Explainf":  func(err error) error { return oops.Explainf(err, "count=%d", 5) },
		"AddCauses": func(err error) error { return oops.AddCauses(err, oops.CauseAuth) },
		"Pathf":     func(err error) error { return oops.Pathf(err, "user/%d", 42) },
	}

	for name, helper := range helpers {
		t.Run(name+"/nil input returns nil", func(t *testing.T) {
			t.Parallel()
			if result := helper(nil); result != nil {
				t.Fatalf("got %#v", result)
			}
			if result := helper(typedNil); result != nil {
				t.Fatalf("typed nil: got %#v", result)
			}
		})

		t.Run(name+"/oops input is annotated in place", func(t *testing.T) {
			t.Parallel()
			err := oops.Define("test").Yeet()
			got, ok := oops.Native(helper(err))
			if !ok || got != err {
				t.Fatal("expected the same *oops.Error back")
			}
		})

		t.Run(name+"/foreign input is wrapped with ErrForeign", func(t *testing.T) {
			t.Parallel()
			original := errors.New("plain")
			result := helper(original)
			got, ok := oops.Native(result)
			if !ok || got.Definition() != oops.ErrForeign {
				t.Fatalf("expected ErrForeign wrapper, got %v", result)
			}
			if !errors.Is(result, original) {
				t.Fatal("original should be reachable")
			}
		})
	}

	t.Run("Explainf on foreign input", func(t *testing.T) {
		t.Parallel()
		result := oops.Explainf(errors.New("plain"), "count=%d", 5)
		if result.Error() != "foreign: count=5" {
			t.Fatalf("got %q", result.Error())
		}
	})

	t.Run("AddCauses keeps existing causes", func(t *testing.T) {
		t.Parallel()
		result, _ := oops.Native(oops.AddCauses(errors.New("plain"), oops.CauseAuth))
		if got := result.Causes(); !slices.Equal(got, []string{oops.CauseInternal, oops.CauseAuth}) {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("Pathf sets the path", func(t *testing.T) {
		t.Parallel()
		result, _ := oops.Native(oops.Pathf(oops.Define("test").Yeet(), "user/%d", 42))
		if result.Path() != "user/42" {
			t.Fatalf("got %q", result.Path())
		}
	})

	t.Run("Explainf formats explanation", func(t *testing.T) {
		t.Parallel()
		result, _ := oops.Native(oops.Explainf(oops.Define("test").Yeet(), "count=%d", 5))
		if result.Explanation() != "count=5" {
			t.Fatalf("got %q", result.Explanation())
		}
	})
}

func TestAs(t *testing.T) {
	t.Parallel()

	t.Run("nil returns false", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		if _, ok := oops.As(nil, oops.Define("test")); ok {
			t.Fatal("expected false for nil")
		}
		if _, ok := oops.As(typedNil, oops.Define("test")); ok {
			t.Fatal("expected false for typed nil")
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
		if !ok || found != innerErr {
			t.Fatal("expected the inner error")
		}
	})

	t.Run("wrapped tree without the target", func(t *testing.T) {
		t.Parallel()
		outerErr := oops.Define("outer").Wrap(oops.Define("inner").Yeet())
		if found, ok := oops.As(outerErr, oops.Define("unrelated")); ok || found != nil {
			t.Fatalf("expected no match, got %v", found)
		}
	})

	t.Run("via stdlib wrapping", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("target")
		innerErr := def.Yeetf("inner")
		joined := errors.Join(errors.New("other"), innerErr)

		found, ok := oops.As(joined, def)
		if !ok || found != innerErr {
			t.Fatal("expected match through errors.Join")
		}
	})

	t.Run("errors.Join without the target", func(t *testing.T) {
		t.Parallel()
		joined := errors.Join(errors.New("other"), oops.Define("unrelated").Yeet())
		if _, ok := oops.As(joined, oops.Define("target")); ok {
			t.Fatal("expected no match")
		}
	})

	t.Run("skips typed-nil nodes", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		def := oops.Define("target")
		targetErr := def.Yeet()
		joined := errors.Join(typedNil, fmt.Errorf("wrap: %w", typedNil), targetErr)

		found, ok := oops.As(joined, def)
		if !ok || found != targetErr {
			t.Fatal("expected match past typed-nil nodes")
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
		var typedNil *oops.Error
		if oops.Nest(oops.Define("test"), nil, typedNil) != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("creates parent with wrapped children", func(t *testing.T) {
		t.Parallel()
		parent := oops.Define("parent")
		err1 := oops.Define("child1").Yeet()
		err2 := oops.Define("child2").Yeet()

		result, ok := oops.Native(oops.Nest(parent, err1, err2))
		if !ok {
			t.Fatal("expected an *oops.Error")
		}
		if result.Definition() != parent {
			t.Fatalf("expected parent definition, got %q", result.Code())
		}
		if got := result.Unwrap(); len(got) != 2 || got[0] != err1 || got[1] != err2 {
			t.Fatalf("got %v", got)
		}
	})

	t.Run("skips nil errors in list", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		parent := oops.Define("parent")
		err1 := oops.Define("child").Yeet()

		result, _ := oops.Native(oops.Nest(parent, nil, err1, typedNil))
		if got := result.Unwrap(); len(got) != 1 || got[0] != err1 {
			t.Fatalf("got %v", got)
		}
	})
}

func TestPresets(t *testing.T) {
	t.Parallel()

	t.Run("ErrForeign", func(t *testing.T) {
		t.Parallel()
		err := oops.ErrForeign.Yeet()
		if err.Code() != "foreign" {
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
