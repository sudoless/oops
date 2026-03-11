package oops_test

import (
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestError_Path(t *testing.T) {
	t.Parallel()

	t.Run("WithPathf with args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().WithPathf("user/%d", 42)
		if err.Path() != "user/42" {
			t.Fatalf("got %q", err.Path())
		}

		args := err.PathArgs()
		if len(args) != 1 {
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
		err = err.WithPathf("static")
		args := err.PathArgs()
		if args != nil {
			t.Fatalf("expected nil args, got %v", args)
		}
	})
}
