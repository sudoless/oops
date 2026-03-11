package oops_test

import (
	"errors"
	"testing"

	"go.sdls.io/oops/pkg/oops"
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

	t.Run("collects oops errors", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		child := oops.Define("child")
		addf(child.Yeet(), "step1")
		addf(child.Yeetf("detail"), "step2")

		result := finish()
		if result == nil {
			t.Fatal("expected non-nil")
		}

		wrapped := result.Unwrap()
		if len(wrapped) != 2 {
			t.Fatalf("expected 2 wrapped, got %d", len(wrapped))
		}

		// Check that paths were set
		first, ok := wrapped[0].(*oops.Error)
		if !ok {
			t.Fatal("expected *Error")
		}
		if first.Path() != "step1" {
			t.Fatalf("expected path %q, got %q", "step1", first.Path())
		}
	})

	t.Run("skips nil errors", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		addf(nil, "ignored")

		if finish() != nil {
			t.Fatal("expected nil")
		}
	})

	t.Run("wraps non-oops errors", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		addf(errors.New("stdlib error"), "path")

		result := finish()
		if result == nil {
			t.Fatal("expected non-nil")
		}

		wrapped := result.Unwrap()
		if len(wrapped) != 1 {
			t.Fatalf("expected 1 wrapped, got %d", len(wrapped))
		}

		oErr, ok := wrapped[0].(*oops.Error)
		if !ok {
			t.Fatal("expected *Error wrapping stdlib error")
		}
		if !errors.Is(oErr, oops.ErrUncaught) {
			t.Fatal("expected ErrUncaught wrapping")
		}
	})

	t.Run("path with format args", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		child := oops.Define("child")
		addf(child.Yeet(), "item/%d", 42)

		result := finish()
		wrapped := result.Unwrap()
		first, _ := wrapped[0].(*oops.Error)
		if first.Path() != "item/42" {
			t.Fatalf("expected path %q, got %q", "item/42", first.Path())
		}
	})

	t.Run("empty path is not added", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		finish, addf := def.Collect()

		child := oops.Define("child")
		addf(child.Yeet(), "")

		result := finish()
		wrapped := result.Unwrap()
		first, _ := wrapped[0].(*oops.Error)
		if first.Path() != "" {
			t.Fatalf("expected empty path, got %q", first.Path())
		}
	})
}
