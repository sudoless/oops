package oops

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var (
	errTest              = Define(BlameServer, NamespaceTest, ReasonInternal)
	errTestHelp          = Define(BlameServer, NamespaceTest, ReasonLegal, "check article 31.40.m", "help 2", "help 3")
	errTestExplainNested = Define(BlameServer, NamespaceTest, ReasonResourceNotFound)
)

func TestError_Help(t *testing.T) {
	t.Parallel()

	err1 := errTest.Yeet()
	if err1.help != "" {
		t.Fatal("error help message must be empty")
	}

	err2 := errTestHelp.Yeet()
	if err2.help != "check article 31.40.m" {
		t.Fatal("error help message does not match expected")
	}
}

func TestErrorDefined_Yeet(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet()

	if err.Error() != "SERVER.TEST.INTERNAL" {
		t.Fatal("err does not have the right code")
	}

	unwrapErr1 := err.Unwrap()
	unwrapErr2 := errors.Unwrap(err)

	if !(unwrapErr1 == nil && unwrapErr2 == nil) {
		t.Fatal("unwrapped errors must be nil")
	}

	if err.Error() != err.Code() {
		t.Fatal("error message does not match error code")
	}

	if err.Error() != err.String() {
		t.Fatal("error message does not match error string")
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

		if err.Error() != "SERVER.TEST.INTERNAL" {
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

func TestExplain(t *testing.T) {
	t.Parallel()

	t.Run("explain nil err", func(t *testing.T) {
		err := Explain(nil, "foo bar baz")
		if err != nil {
			t.Fatal("explain must not create error from nil")
		}
	})

	t.Run("explain new error", func(t *testing.T) {
		err := errors.New("fiz biz")
		errExplained1 := Explain(err, "foo bar")
		errExplained2 := Explain(err, "bar foo")

		if errExplained1.Code() != "UNKNOWN.UNKNOWN.UNEXPECTED" {
			t.Fatal("explained non *Error errors must be UNEXPECTED")
		}

		if !errors.Is(errExplained1, ErrUnexpected) {
			t.Fatal("explained error must not lose inheritance/link to errorDefined")
		}

		if errors.Is(errExplained1, errTest) {
			t.Fatal("explained error must not have unrelated inheritance/link")
		}

		if !errors.Is(errExplained1, errExplained2) {
			t.Fatal("explained error must have sibling inheritance/link")
		}
	})
}

func TestError_MarshalJSON(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		errEmpty := &Error{}
		j, err := errEmpty.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"UNKNOWN.UNKNOWN.UNKNOWN\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("yeet simple", func(t *testing.T) {
		errYeet := errTest.Yeet()
		j, err := errYeet.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"SERVER.TEST.INTERNAL\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("yeet help", func(t *testing.T) {
		errYeet := errTestHelp.Yeet()
		j, err := errYeet.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"SERVER.TEST.LEGAL\",\"help\":\"check article 31.40.m\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("wrap", func(t *testing.T) {
		errWrap := errTest.Wrap(errors.New("foo bar"))

		j, err := errWrap.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"SERVER.TEST.INTERNAL\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("explain", func(t *testing.T) {
		errExplained := errTest.YeetExplain("explain 1")

		j, err := errExplained.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"SERVER.TEST.INTERNAL\",\"explain\":\"explain 1\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("explain x3", func(t *testing.T) {
		errExplained := errTest.YeetExplain("explain 1")
		_ = Explain(errExplained, "explain 2")
		_ = Explain(errExplained, "explain 3")

		j, err := errExplained.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"SERVER.TEST.INTERNAL\",\"explain\":\"explain 1, explain 2, explain 3\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("wrap explain empty", func(t *testing.T) {
		errExplained := errTest.WrapExplain(errors.New("foo bar"), "explain 1")

		j, err := errExplained.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"SERVER.TEST.INTERNAL\",\"explain\":\"explain 1\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})
}

func TestError_Is(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		err := errTest.Yeet()

		if errors.Is(err, nil) {
			t.Fatal("error must not match Is nil")
		}

		if err.Is(nil) {
			t.Fatal("error must not match Is nil")
		}
	})

	t.Run("is nil", func(t *testing.T) {
		var err *Error

		if !err.Is(nil) {
			t.Fatal("error in this case must match Is nil")
		}
	})

	t.Run("is parent", func(t *testing.T) {
		parent := errors.New("foo bar")
		err := errTest.Wrap(parent)

		if !errors.Is(err, parent) {
			t.Fatal("error must match parent on Is")
		}

		if !err.Is(parent) {
			t.Fatal("error must match parent on Is")
		}
	})

	t.Run("nil parent", func(t *testing.T) {
		parent := errors.New("foo bar")
		err := errTest.Wrap(nil)

		if errors.Is(err, parent) {
			t.Fatal("error must not match parent on Is")
		}

		if err.Is(parent) {
			t.Fatal("error must not match parent on Is")
		}
	})
}

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

func TestError_Multi(t *testing.T) {
	t.Parallel()

	t.Run("x3 + x2", func(t *testing.T) {
		errYeet := errTest.Yeet().Multi("reason 1", "reason 2", "reason 3")
		_ = errYeet.Multi("reason 4", "reason 5")

		if len(errYeet.multi) != 5 {
			t.Fatal("error multi length not ok")
		}

		if !reflect.DeepEqual([]string{"reason 1", "reason 2", "reason 3", "reason 4", "reason 5"}, errYeet.multi) {
			t.Fatal("error multi does not match string array")
		}
	})

	t.Run("none", func(t *testing.T) {
		errYeet := errTest.Yeet().Multi()
		if errYeet.multi != nil {
			t.Fatal("error multi must be nil")
		}
	})

	t.Run("json", func(t *testing.T) {
		errYeet := errTest.Yeet().Multi("foo", "bar", "baz")

		j, err := errYeet.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"SERVER.TEST.INTERNAL\",\"multi\":[\"foo\",\"bar\",\"baz\"]}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})
}

func Test_stack(t *testing.T) {
	t.Run("no stack", func(t *testing.T) {
		err := errTest.NoStack().Yeet()
		if err.trace != nil {
			t.Fatal("no stack error should have no trace")
		}
	})

	t.Run("normal depth", func(t *testing.T) {
		err := errTest.Yeet()
		if err.trace == nil {
			t.Fatal("error should have stack trace")
		}

		for _, trace := range err.trace {
			fmt.Println(trace)
		}
	})
}

func Test_Explain_nested(t *testing.T) {
	t.Parallel()

	err := testExplainCaller()
	if err == nil {
		t.Fatal("expected non nil error")
	}

	if !errors.Is(err, errTestExplainNested) {
		t.Fatal("expected error to be errTestExplainNested")
	}

	oopsErr, ok := err.(*Error)
	if !ok {
		t.Fatal("expected *oops.Error")
	}

	if expln := oopsErr.Explanation(); expln != "source not found, middleware 2 applied, midd 1 happened, performing middle 0 action, caller explaining" {
		t.Fatal("wrong error explanation", expln)
	}
}

func testExplainCaller() error {
	err := testExplainMiddle0()
	if err != nil {
		return Explain(err, "caller explaining")
	}
	return nil
}

func testExplainMiddle0() error {
	err := testExplainMiddle1()
	if err != nil {
		return Explain(err, "performing middle 0 action")
	}
	return nil
}

func testExplainMiddle1() error {
	err := testExplainMiddle2()
	if err != nil {
		return Explain(err, "midd 1 happened")
	}
	return nil
}

func testExplainMiddle2() error {
	err := testExplainSource()
	if err != nil {
		return Explain(err, "middleware 2 applied")
	}
	return nil
}

func testExplainSource() error {
	return errTestExplainNested.YeetExplain("source not found")
}

func Test_String(t *testing.T) {
	t.Parallel()

	t.Run("yeet explain", func(t *testing.T) {
		err := errTest.YeetExplain("foo bar")
		if String(err) != "SERVER.TEST.INTERNAL foo bar" {
			t.Fatal("unexpected error string from oops.String:", String(err))
		}
	})

	t.Run("nil", func(t *testing.T) {
		if s := String(nil); s != "" {
			t.Fatal("oops.String on nil should be empty not:", s)
		}
	})

	t.Run("builtin error", func(t *testing.T) {
		err := errors.New("foo bar")
		if String(err) != "foo bar" {
			t.Fatal("unexpected error string from oops.String:", String(err))
		}
	})
}

func Test_ExplainFmt(t *testing.T) {
	t.Parallel()

	t.Run("yeet fmt", func(t *testing.T) {
		err := errTest.YeetExplainFmt("foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})

	t.Run("wrap fmt", func(t *testing.T) {
		err := errTest.WrapExplainFmt(errors.New("fiz"), "foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})
}
