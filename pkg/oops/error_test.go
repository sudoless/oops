package oops_test

import (
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

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
