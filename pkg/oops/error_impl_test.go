package oops_test

import (
	"errors"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func Test_stack(t *testing.T) {
	t.Parallel()

	t.Run("no stack", func(t *testing.T) {
		t.Parallel()

		err := errTest.Yeet()
		if err.Trace() != nil {
			t.Fatal("no stack error should have no trace")
		}
	})

	t.Run("normal depth", func(t *testing.T) {
		t.Parallel()

		err := errTestTrace.Yeet()
		if err.Trace() == nil {
			t.Fatal("error should have stack trace")
		}

		for _, trace := range err.Trace() {
			t.Log(trace)
		}
	})
}

func Test_ExplainFmt(t *testing.T) {
	t.Parallel()

	t.Run("yeet fmt", func(t *testing.T) {
		t.Parallel()

		err := errTest.Yeetf("foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})

	t.Run("wrap fmt", func(t *testing.T) {
		t.Parallel()

		err := errTest.Wrapf(errors.New("fiz"), "foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})

	t.Run("explain", func(t *testing.T) {
		t.Parallel()

		err := errors.New("new")
		out := oops.Explainf(err, "foo %s", "bar")
		msg := out.Error()

		if msg != "uncaught unwrapped: foo bar" {
			t.Fatalf("unexpected error message('%s')", msg)
		}
	})
}

func TestError_String(t *testing.T) {
	t.Parallel()

	err := errTest.Yeetf("foobar")
	_ = oops.Explainf(err, "fiz")
	_ = oops.Explainf(err, "fuz")

	s := err.Error()
	if s != "foobar, fiz, fuz" {
		t.Fatalf("error string does not match expectations('%s')", s)
	}
}

func TestError_Explain(t *testing.T) {
	t.Parallel()

	err := errTest.Yeetf("foo bar")
	err.Explainf("baz")

	if err.Explanation() != "foo bar, baz" {
		t.Errorf("expected explanation 'foo bar, baz', got %s", err.Explanation())
	}

	err.Explainf("id=%d", 123)
	if err.Explanation() != "foo bar, baz, id=123" {
		t.Errorf("expected explanation 'foo bar, baz, id=123', got %s", err.Explanation())
	}
}

func BenchmarkError_String(b *testing.B) {
	b.ReportAllocs()

	err := errTest.Yeet()

	for iter := 0; iter <= b.N; iter++ {
		_ = err.Error()
	}
}

func TestError_Set(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet()
	_ = err.Set("request_id", "abc-123")

	v, ok := err.Get("request_id")
	if !ok {
		t.Fatal("Set key must be retrievable with Get")
	}
	if v.(string) != "abc-123" {
		t.Fatalf("unexpected value: %v", v)
	}

	// overwrite
	_ = err.Set("request_id", "xyz-789")
	v, ok = err.Get("request_id")
	if !ok {
		t.Fatal("overwritten key must still be present")
	}
	if v.(string) != "xyz-789" {
		t.Fatalf("unexpected overwritten value: %v", v)
	}
}

func TestError_Get_missing(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet()
	_, ok := err.Get("nonexistent")
	if ok {
		t.Fatal("Get must return false for missing key")
	}
}

func TestError_Append(t *testing.T) {
	t.Parallel()

	parent := errTest.Yeet()
	c1 := oops.Define().Yeetf("child one")
	c2 := oops.Define().Yeetf("child two")

	result := parent.Append(c1, c2)
	if result != parent {
		t.Fatal("Append must return the same error")
	}

	nested := parent.Nested()
	if len(nested) != 2 {
		t.Fatalf("expected 2 nested errors, got %d", len(nested))
	}
	if nested[0].Explanation() != "child one" {
		t.Fatalf("unexpected nested[0] explanation: %q", nested[0].Explanation())
	}
	if nested[1].Explanation() != "child two" {
		t.Fatalf("unexpected nested[1] explanation: %q", nested[1].Explanation())
	}
}

func TestError_Path(t *testing.T) {
	t.Parallel()

	t.Run("with args", func(t *testing.T) {
		t.Parallel()

		err := errTest.Yeet()
		_ = err.PathSetf("user/%d/profile", 42)

		if err.Path() != "user/42/profile" {
			t.Fatalf("unexpected Path: %q", err.Path())
		}
		args := err.PathArgs()
		if len(args) != 1 || args[0].(int) != 42 {
			t.Fatalf("unexpected PathArgs: %v", args)
		}
	})

	t.Run("without args", func(t *testing.T) {
		t.Parallel()

		err := errTest.Yeet()
		_ = err.PathSetf("static/path")

		if err.Path() != "static/path" {
			t.Fatalf("unexpected Path: %q", err.Path())
		}
		if err.PathArgs() != nil {
			t.Fatalf("PathArgs must be nil when no args given: %v", err.PathArgs())
		}
	})
}

func TestError_Explainf_emptyFormat(t *testing.T) {
	t.Parallel()

	err := errTest.Yeetf("initial")
	err.Explainf("")

	if err.Explanation() != "initial" {
		t.Fatalf("empty Explainf must not modify explanation, got %q", err.Explanation())
	}
}

func BenchmarkError_wrapExplain(b *testing.B) {
	b.ReportAllocs()

	originalErr := errors.New("original error")

	b.ResetTimer()
	for iter := 0; iter <= b.N; iter++ {
		_ = benchmarkNested1(originalErr)
	}
}

func benchmarkNested1(original error) error {
	if err := benchmarkNested2(original); err != nil {
		return oops.Explainf(err, "nested error 1")
	}

	return nil
}

func benchmarkNested2(original error) error {
	if err := benchmarkNested3(original); err != nil {
		return oops.Explainf(err, "nested error 2")
	}

	return nil
}

func benchmarkNested3(original error) error {
	err := benchmarkNested4(original)
	if err != nil {
		return oops.Explainf(err, "nested error 3")
	}

	return nil
}

func benchmarkNested4(original error) error {
	return errTestBenchmark.Wrapf(original, "benchmarkNested4 returned wrapped original error")
}
