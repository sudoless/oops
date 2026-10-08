package oops_test

import (
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

		args := err.PathArgs()
		if len(args) != 1 {
			t.Fatalf("expected 1 path arg, got %v", args)
		}
	})

	t.Run("Pathf static sets nil args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet().Pathf("static")
		args := err.PathArgs()
		if args != nil {
			t.Fatalf("expected nil args, got %v", args)
		}
	})

	t.Run("Pathf no args sets nil args", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		err = err.Pathf("static")
		args := err.PathArgs()
		if args != nil {
			t.Fatalf("expected nil args, got %v", args)
		}
	})
}
