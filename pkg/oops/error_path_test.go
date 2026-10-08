package oops_test

import (
	"slices"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

func TestError_Path(t *testing.T) {
	t.Parallel()

	t.Run("Pathf with args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Pathf("user/%d", 42)
		if err.Path() != "user/42" {
			t.Fatalf("got %q", err.Path())
		}
		if args := err.PathArgs(); !slices.Equal(args, []any{42}) {
			t.Fatalf("expected path args [42], got %v", args)
		}
	})

	t.Run("Pathf no args sets nil args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Pathf("static")
		if err.Path() != "static" {
			t.Fatalf("got %q", err.Path())
		}
		if args := err.PathArgs(); args != nil {
			t.Fatalf("expected nil args, got %v", args)
		}
	})

	t.Run("Pathf overwrites path and args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Pathf("user/%d", 42).Pathf("static")
		if err.Path() != "static" {
			t.Fatalf("got %q", err.Path())
		}
		if args := err.PathArgs(); args != nil {
			t.Fatalf("expected stale args to be cleared, got %v", args)
		}
	})

	t.Run("empty format clears the path", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Pathf("user/%d", 42).Pathf("")
		if err.Path() != "" || err.PathArgs() != nil {
			t.Fatalf("got %q %v", err.Path(), err.PathArgs())
		}
	})

	t.Run("escaped percent without args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Pathf("100%%")
		if err.Path() != "100%" {
			t.Fatalf("got %q", err.Path())
		}
	})

	t.Run("Pathf keeps its own copy of the args", func(t *testing.T) {
		t.Parallel()
		args := []any{1}
		err := oops.Define("test").Yeet().Pathf("item/%d", args...)
		args[0] = 2
		if got := err.PathArgs(); len(got) != 1 || got[0] != 1 {
			t.Fatalf("got %v", got)
		}
	})
}
