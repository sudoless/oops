package oops_test

import (
	"errors"
	"fmt"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

var (
	errTest              = oops.Define("code", "test.err_test")
	errTestExplainNested = oops.Define("code", "test.err_test_explain_nested")
	errTestTrace         = oops.Define("code", "test.err_test_trace").Trace()
	errTestBenchmark     = oops.Define("code", "test.err_test_benchmark", "status", 418)
)

func TestErrorDefined_Error(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("define Error method did not panic")
		}
	}()

	str := errTest.Error()

	if str != "" {
		t.Fatal("define Error method managed to return a string")
	}
}

func TestErrorDefined_customFormatter(t *testing.T) {
	t.Parallel()

	t.Run("simple", func(t *testing.T) {
		customFormatter := func(err oops.Error) string {
			return "custom formatter"
		}

		errd := oops.Define("code", "test.custom_formatter").Formatter(customFormatter)

		err := errd.Yeet()
		if err.Error() != "custom formatter" {
			t.Fatal("expected custom formatter to be used")
		}
	})

	t.Run("type and code", func(t *testing.T) {
		customFormatter := func(err oops.Error) string {
			m := err.GetAll()
			code := m["code"].(string)
			typ := m["type"].(string)

			explain := err.Explanation()
			if explain == "" {
				return fmt.Sprintf("[%s] %s", typ, code)
			}

			return fmt.Sprintf("[%s] %s : %s", typ, code, explain)
		}

		errd := oops.Define("type", "test", "code", "custom_formatter").Formatter(customFormatter)

		err := errd.Yeet()
		if err.Error() != "[test] custom_formatter" {
			t.Fatal("expected custom formatter to be used")
		}

		err = oops.Explainf(err, "foobar")
		if err.Error() != "[test] custom_formatter : foobar" {
			t.Fatal("expected custom formatter to be used")
		}

		err = oops.Explainf(err, "fiz")
		if err.Error() != "[test] custom_formatter : foobar, fiz" {
			t.Fatal("expected custom formatter to be used")
		}

		errw := errd.Wrapf(errors.New("fiz"), "foobar")
		if errw.Error() != "[test] custom_formatter : foobar" {
			t.Fatal("expected custom formatter to be used")
		}
	})

	t.Run("unwrap parent", func(t *testing.T) {
		customFormatter := func(err oops.Error) string {
			return "custom formatter: " + err.Unwrap().Error()
		}

		errd := oops.Define("code", "test.custom_formatter").Formatter(customFormatter)

		errw := errd.Wrapf(errors.New("fiz"), "foobar")
		if errw.Error() != "custom formatter: fiz" {
			t.Fatal("expected custom formatter to be used")
		}
	})
}

func TestErrorDefined_Yeet(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet()

	if code, ok := err.Get("code"); !ok || code != "test.err_test" {
		t.Fatal("err does not have the right code")
	}

	unwrapErr1 := err.Unwrap()
	unwrapErr2 := errors.Unwrap(err)

	if !(unwrapErr1 == nil && unwrapErr2 == nil) {
		t.Fatal("unwrapped errors must be nil")
	}

	if err.Explanation() != "" {
		t.Fatal("explanation must be empty")
	}
}

func TestErrorDefined_Wrap(t *testing.T) {
	t.Parallel()

	t.Run("new error", func(t *testing.T) {
		err := errTest.Wrap(errors.New("failed dial target host"))
		if err == nil {
			t.Fatal("err cannot be nil after wrap")
		}

		if code, ok := err.Get("code"); !ok || code != "test.err_test" {
			t.Fatal("err does not have the right code")
		}

		unwrapErr1 := err.Unwrap()
		unwrapErr2 := errors.Unwrap(err)

		if unwrapErr1 == nil || unwrapErr2 == nil {
			t.Fatal("unwrapped errors are nil")
		}

		if unwrapErr1 != unwrapErr2 {
			t.Fatal("unwrapped errors are not equal")
		}

		if unwrapErr1.Error() != "failed dial target host" {
			t.Fatal("unwrapped error message does not match expected")
		}

		if err.Explanation() != "" {
			t.Fatal("explanation must be empty")
		}
	})

	t.Run("errors is", func(t *testing.T) {
		parent := errors.New("daddy error")
		err := errTest.Wrap(parent)

		if !errors.Is(err, parent) {
			t.Fatal("errors.Is did not match error to parent")
		}
	})
}

func TestErrorDefined_format(t *testing.T) {
	t.Parallel()

	t.Run("yeet fmt", func(t *testing.T) {
		err := errTest.Yeetf("foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})

	t.Run("wrap fmt", func(t *testing.T) {
		err := errTest.Wrapf(errors.New("fiz"), "foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})
}

func TestErrorDefined_Collect(t *testing.T) {
	t.Parallel()

	var (
		errEven = oops.Define()
		errOdd  = oops.Define()
		errNone = oops.Define()
	)

	finish, addf := errTest.Collect()
	for idx := 0; idx < 10; idx++ {
		if idx%2 == 0 {
			addf(errEven, "index %d", idx)
		} else {
			addf(errOdd, "index %d", idx)
		}
	}

	_ = errNone

	err := finish()
	if err.Source() != errTest {
		t.Fatal("err source does not match")
	}

	_, ok := oops.As(err, errEven)
	if ok {
		t.Fatal("oops.As must not check nested errors")
	}

	_, ok = oops.NestedAs(err, errEven)
	if !ok {
		t.Fatal("oops.NestedAs must check nested errors")
	}

	_, ok = oops.NestedAs(err, errOdd)
	if !ok {
		t.Fatal("oops.NestedAs must check nested errors")
	}
}

func TestErrorDefined_Collect_none(t *testing.T) {
	t.Parallel()

	t.Run("zero", func(t *testing.T) {
		t.Parallel()

		finish, _ := errTest.Collect()
		err := finish()
		if err != nil {
			t.Fatal("err must be nil")
		}
	})

	t.Run("all nil", func(t *testing.T) {
		t.Parallel()

		finish, addf := errTest.Collect()
		for idx := 0; idx < 10; idx++ {
			addf(nil, "index %d", idx)
		}

		err := finish()
		if err != nil {
			t.Fatal("err must be nil")
		}
	})
}

func TestErrorDefined_Is(t *testing.T) {
	t.Parallel()

	someErr := errors.New("some err")
	if errTest.Is(someErr) {
		t.Fatal("errTest.Is(someErr)")
	}

	oerr := errTest.Yeet()
	if !errTest.Is(oerr) {
		t.Fatal("expected errTest.Is(oerr)")
	}

	if errTest.Is(nil) {
		t.Fatal("errTest.Is(nil)")
	}
}
