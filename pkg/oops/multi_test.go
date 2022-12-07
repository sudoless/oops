package oops

import "testing"

func TestMulti(t *testing.T) {
	t.Parallel()

	multi := Multi(
		errTest.Yeet("test"),
		errTestHelp.Yeet("test help"),
		errTestTrace.Yeet("test trace"),
	)

	errs := multi.(*multipleErrors).Unwrap()
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(errs))
	}

	for _, err := range errs {
		t.Log(err.Error())
	}
}
