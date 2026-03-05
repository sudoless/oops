package oops_test

import (
	"errors"
	"fmt"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestExplain_nested(t *testing.T) {
	t.Parallel()

	testExplainSource := func() error {
		return errTestExplainNested.Yeetf("source not found")
	}

	testExplainMiddle2 := func() error {
		err := testExplainSource()
		if err != nil {
			return oops.Explainf(err, "middleware 2 applied")
		}
		return nil
	}

	testExplainMiddle1 := func() error {
		err := testExplainMiddle2()
		if err != nil {
			return oops.Explainf(err, "midd 1 happened")
		}
		return nil
	}

	testExplainMiddle0 := func() error {
		err := testExplainMiddle1()
		if err != nil {
			return oops.Explainf(err, "performing middle 0 action")
		}
		return nil
	}

	testExplainCaller := func() error {
		err := testExplainMiddle0()
		if err != nil {
			return oops.Explainf(err, "caller explaining")
		}
		return nil
	}

	err := testExplainCaller()
	if err == nil {
		t.Fatal("expected non nil error")
	}

	if !errors.Is(err, errTestExplainNested) {
		t.Fatal("expected error to be errTestExplainNested")
	}

	oopsErr, ok := err.(oops.Error) //nolint:errorlint
	if !ok {
		t.Fatal("expected *oops.Error")
	}

	if expln := oopsErr.Explanation(); expln != "source not found, middleware 2 applied, midd 1 happened, performing middle 0 action, caller explaining" {
		t.Fatal("wrong error explanation", expln)
	}
}

func TestExplain(t *testing.T) {
	t.Parallel()

	t.Run("explain nil err", func(t *testing.T) {
		t.Parallel()

		err := oops.Explainf(nil, "foo bar baz")
		if err != nil {
			t.Fatal("explain must not create error from nil")
		}
	})

	t.Run("explain nil *oops.Error", func(t *testing.T) {
		t.Parallel()

		returnNil := func() oops.Error {
			return nil
		}

		err := oops.Explainf(returnNil(), "foo bar baz")
		if err != nil {
			t.Fatal("explain must not create error from nil")
		}
	})

	t.Run("explain new error", func(t *testing.T) {
		t.Parallel()

		err := errors.New("fiz biz")
		errExplained1 := oops.Explainf(err, "foo bar")
		errExplained2 := oops.Explainf(err, "bar foo")

		if !errors.Is(errExplained1, oops.ErrUncaught) {
			t.Fatal("explained error must not lose inheritance/link to ErrorDefined")
		}

		if errors.Is(errExplained1, errTest) {
			t.Fatal("explained error must not have unrelated inheritance/link")
		}

		if !errors.Is(errExplained1, errExplained2) {
			t.Fatal("explained error must have sibling inheritance/link")
		}
	})

	t.Run("format", func(t *testing.T) {
		t.Parallel()

		err := errors.New("new")
		out := oops.Explainf(err, "foo %s", "bar")
		msg := out.Error()

		if msg != "uncaught unwrapped: foo bar" {
			t.Fatalf("unexpected error message('%s')", msg)
		}
	})
}

type minimalError struct{}

func (m minimalError) Error() string { return "minimal error" }

func TestAs_minimal(t *testing.T) {
	t.Parallel()

	ohno := func() error {
		return minimalError{}
	}

	err := ohno()

	oerr, ok := oops.As(err, errTest)
	if ok {
		t.Fatal("error cannot be errTest")
	}

	t.Log(oerr)
}

func TestAs_fmtErrorfWrap(t *testing.T) {
	t.Parallel()

	inner := errTest.Yeetf("inner explanation")
	wrapped := fmt.Errorf("context: %w", inner)

	got, ok := oops.As(wrapped, errTest)
	if !ok {
		t.Fatal("oops.As must traverse fmt.Errorf %w wrapper")
	}
	if got.Explanation() != "inner explanation" {
		t.Fatalf("unexpected explanation: %q", got.Explanation())
	}
}

func TestAs_fmtErrorfWrap_notFound(t *testing.T) {
	t.Parallel()

	other := oops.Define("code", "other")
	inner := errTest.Yeetf("inner")
	wrapped := fmt.Errorf("context: %w", inner)

	_, ok := oops.As(wrapped, other)
	if ok {
		t.Fatal("oops.As must not find unrelated defined error through wrapper")
	}
}

func TestAssertAny(t *testing.T) {
	t.Parallel()

	t.Run("oops error", func(t *testing.T) {
		t.Parallel()

		err := errTest.Yeetf("hello")
		got, ok := oops.AssertAny(err)
		if !ok {
			t.Fatal("AssertAny must succeed for oops.Error")
		}
		if got.Explanation() != "hello" {
			t.Fatalf("unexpected explanation: %q", got.Explanation())
		}
	})

	t.Run("non-oops error", func(t *testing.T) {
		t.Parallel()

		_, ok := oops.AssertAny(errors.New("plain"))
		if ok {
			t.Fatal("AssertAny must fail for non-oops error")
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		got, ok := oops.AssertAny(nil)
		if ok {
			t.Fatal("AssertAny must fail for nil")
		}
		if got != nil {
			t.Fatal("AssertAny must return nil Error for nil input")
		}
	})
}

func TestMustAny(t *testing.T) {
	t.Parallel()

	t.Run("oops error", func(t *testing.T) {
		t.Parallel()

		err := errTest.Yeetf("hello")
		got := oops.MustAny(err)
		if got == nil {
			t.Fatal("MustAny must return non-nil for oops.Error")
		}
		if got.Explanation() != "hello" {
			t.Fatalf("unexpected explanation: %q", got.Explanation())
		}
	})

	t.Run("non-oops error wraps with ErrUncaught", func(t *testing.T) {
		t.Parallel()

		plain := errors.New("plain error")
		got := oops.MustAny(plain)
		if got == nil {
			t.Fatal("MustAny must return non-nil for non-oops error")
		}
		if !errors.Is(got, oops.ErrUncaught) {
			t.Fatal("MustAny must wrap non-oops error with ErrUncaught")
		}
		if !errors.Is(got, plain) {
			t.Fatal("MustAny must preserve original error in unwrap chain")
		}
	})

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		got := oops.MustAny(nil)
		if got != nil {
			t.Fatal("MustAny must return nil for nil input")
		}
	})
}

func TestNest(t *testing.T) {
	t.Parallel()

	t.Run("with nested", func(t *testing.T) {
		t.Parallel()

		var (
			errParent = oops.Define("code", "parent")
			errChild1 = oops.Define("code", "child1")
			errChild2 = oops.Define("code", "child2")
		)

		c1 := errChild1.Yeetf("first")
		c2 := errChild2.Yeetf("second")
		parent := oops.Nest(errParent, c1, c2)

		if parent == nil {
			t.Fatal("Nest must return non-nil when nested errors provided")
		}
		if parent.Source() != errParent {
			t.Fatal("Nest source must match provided ErrorDefined")
		}

		_, ok := oops.NestedAs(parent, errChild1)
		if !ok {
			t.Fatal("NestedAs must find child1")
		}
		_, ok = oops.NestedAs(parent, errChild2)
		if !ok {
			t.Fatal("NestedAs must find child2")
		}
	})

	t.Run("nil source", func(t *testing.T) {
		t.Parallel()

		c := errTest.Yeet()
		got := oops.Nest(nil, c)
		if got != nil {
			t.Fatal("Nest must return nil when source is nil")
		}
	})

	t.Run("no nested", func(t *testing.T) {
		t.Parallel()

		got := oops.Nest(errTest)
		if got != nil {
			t.Fatal("Nest must return nil when no nested errors provided")
		}
	})
}

func TestNestedIs(t *testing.T) {
	t.Parallel()

	var (
		errParent = oops.Define("code", "parent")
		errChild  = oops.Define("code", "child")
		errOther  = oops.Define("code", "other")
	)

	t.Run("found", func(t *testing.T) {
		t.Parallel()

		parent := oops.Nest(errParent, errChild.Yeet())
		if !oops.NestedIs(parent, errChild) {
			t.Fatal("NestedIs must find child")
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		parent := oops.Nest(errParent, errChild.Yeet())
		if oops.NestedIs(parent, errOther) {
			t.Fatal("NestedIs must not find unrelated error")
		}
	})

	t.Run("nil err nil target", func(t *testing.T) {
		t.Parallel()

		if !oops.NestedIs(nil, nil) {
			t.Fatal("NestedIs(nil, nil) must return true")
		}
	})

	t.Run("nil err non-nil target", func(t *testing.T) {
		t.Parallel()

		if oops.NestedIs(nil, errChild) {
			t.Fatal("NestedIs(nil, nonNil) must return false")
		}
	})

	t.Run("non-oops error unwraps to find nested", func(t *testing.T) {
		t.Parallel()

		inner := oops.Nest(errParent, errChild.Yeet())
		wrapped := fmt.Errorf("wrap: %w", inner)
		if !oops.NestedIs(wrapped, errChild) {
			t.Fatal("NestedIs must traverse unwrap chain to find nested")
		}
	})

	t.Run("nil oops error", func(t *testing.T) {
		t.Parallel()

		var nilOops oops.Error
		if !oops.NestedIs(nilOops, nil) {
			t.Fatal("NestedIs(nilOops, nil) must return true")
		}
	})
}
