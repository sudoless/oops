package oops

import "testing"

func TestMulti(t *testing.T) {
	t.Parallel()

	multi := Join(
		errTest.Yeet("test"),
		errTestHelp.Yeet("test help"),
		errTestTrace.Yeet("test trace"),
	)

	errs := multi.Errs()
	if len(errs) != 3 {
		t.Fatalf("expected 3 errors, got %d", len(errs))
	}

	for _, err := range errs {
		t.Log(err.Error())
	}
}
