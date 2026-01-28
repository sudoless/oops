package oops_test

import (
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestErrorDefined_Set(t *testing.T) {
	t.Parallel()

	err := oops.Define().Set("foo", "bar")
	yerr := err.Yeet()

	v, ok := yerr.Get("foo")
	if !ok {
		t.Fatalf("foo was not defined")
	}
	if v.(string) != "bar" {
		t.Fatalf("foo was not %v, got %v", `"bar"`, v.(string))
	}
}
