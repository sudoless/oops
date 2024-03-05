package oops_test

import (
	"errors"
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

	oopsErr, ok := err.(oops.Error)
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
		err := oops.Explainf(nil, "foo bar baz")
		if err != nil {
			t.Fatal("explain must not create error from nil")
		}
	})

	t.Run("explain nil *oops.Error", func(t *testing.T) {
		returnNil := func() oops.Error {
			return nil
		}

		err := oops.Explainf(returnNil(), "foo bar baz")
		if err != nil {
			t.Fatal("explain must not create error from nil")
		}
	})

	t.Run("explain new error", func(t *testing.T) {
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
		err := errors.New("new")
		out := oops.Explainf(err, "foo %s", "bar")
		msg := out.Error()

		if msg != "uncaught unwrapped: foo bar" {
			t.Fatalf("unexpected error message('%s')", msg)
		}
	})
}
