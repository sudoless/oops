package oops_test

import (
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestDefine_defaultFormatter(t *testing.T) {
	t.Parallel()

	want := "oops.Error"
	if got := errTest.Yeet().Error(); got != want {
		t.Fatalf("errTest.Yeet() = %v, want %v", got, want)
	}

	want = "hello world"
	if got := errTest.Yeetf("hello world").Error(); got != want {
		t.Fatalf("errTest.Yeet() = %v, want %v", got, want)
	}

	want = "oops.Error(nil)"
	if got := oops.NilErr.Error(); got != want {
		t.Fatalf("oops.NilErr.Error() = %v, want %v", got, want)
	}
}
