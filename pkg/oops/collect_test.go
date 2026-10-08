package oops_test

import (
	"errors"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

func TestCollect(t *testing.T) {
	t.Parallel()

	t.Run("no errors returns nil", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, _ := def.Collect()
		if finish() != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("collects oops and foreign errors in order", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		child := oops.Define("child")
		original := errors.New("stdlib error")
		addf(child.Yeet(), "step1")
		addf(child.Yeetf("detail"), "item/%d", 2)
		addf(original, "step3")

		result := finish()
		if !errors.Is(result, def) {
			t.Fatal("result must come from the collecting definition")
		}
		if !errors.Is(result, oops.ErrForeign) {
			t.Fatal("foreign child should be reachable as ErrForeign")
		}
		if !errors.Is(result, original) {
			t.Fatal("original foreign error should be reachable")
		}

		parent := mustNative(t, result)
		if parent.Definition() != def {
			t.Fatalf("result definition = %q", parent.Code())
		}

		wrapped := parent.Unwrap()
		if len(wrapped) != 3 {
			t.Fatalf("expected 3 wrapped, got %d", len(wrapped))
		}

		wantPaths := []string{"step1", "item/2", "step3"}
		for idx, want := range wantPaths {
			got, ok := oops.Native(wrapped[idx])
			if !ok {
				t.Fatalf("wrapped[%d] is not an *oops.Error", idx)
			}
			if got.Path() != want {
				t.Errorf("wrapped[%d].Path() = %q, want %q", idx, got.Path(), want)
			}
		}

		foreign := mustNative(t, wrapped[2])
		if foreign.Definition() != oops.ErrForeign {
			t.Fatalf("foreign wrapper definition = %q", foreign.Code())
		}
		if inner := foreign.Unwrap(); len(inner) != 1 || inner[0] != original {
			t.Fatalf("foreign wrapper children = %v", inner)
		}
	})

	t.Run("skips nil errors", func(t *testing.T) {
		t.Parallel()
		var typedNil *oops.Error
		def := oops.Define("test")
		finish, addf := def.Collect()

		addf(nil, "ignored")
		addf(typedNil, "ignored")

		if finish() != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("empty path keeps the existing path", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		addf(oops.Define("child").Yeet().Pathf("original"), "")

		parent := mustNative(t, finish())
		first := mustNative(t, parent.Unwrap()[0])
		if first.Path() != "original" {
			t.Fatalf("expected path %q, got %q", "original", first.Path())
		}
	})

	t.Run("non-empty path overwrites the existing path", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		addf(oops.Define("child").Yeet().Pathf("original"), "replaced")

		parent := mustNative(t, finish())
		first := mustNative(t, parent.Unwrap()[0])
		if first.Path() != "replaced" {
			t.Fatalf("expected path %q, got %q", "replaced", first.Path())
		}
	})

	t.Run("finish result is independent of later adds", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		addf(def.Yeet(), "a")
		first := mustNative(t, finish())
		nested := def.Yeet()
		_ = first.Nest(nested)
		addf(def.Yeet(), "b")

		if got := first.Unwrap(); len(got) != 2 || got[1] != nested {
			t.Fatalf("first result changed by a later add: %v", got)
		}
	})

	t.Run("finish results do not share children storage", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		child := def.Yeet()
		addf(child, "a")
		first := mustNative(t, finish())
		first.Unwrap()[0] = def.Yeet()

		second := mustNative(t, finish())
		if got := second.Unwrap(); len(got) != 1 || got[0] != child {
			t.Fatalf("second result picked up a change to the first: %v", got)
		}
	})
}
