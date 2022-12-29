package oops

import (
	"database/sql"
	"errors"
	"fmt"
	"testing"
)

var (
	errTest              = Define().Code("err_test").Type("test")
	errTestHelp          = Define().Code("err_test_help").Type("test").Help("check article 31.40.m")
	errTestExplainNested = Define().Code("err_test_explain_nested").Type("test")
	errTestTrace         = Define().Code("err_test_trace").Trace().Type("test")
	errTestStatusCode    = Define().Code("err_test_status_code").StatusCode(418).Type("test")
	errTestBenchmark     = Define().Code("err_test_benchmark").Type("benchmark").StatusCode(418)
)

func Test__use(t *testing.T) {
	_ = ErrTODO
}

func testReturnNilOopsError() *Error {
	return nil
}

func TestError_Help(t *testing.T) {
	t.Parallel()

	err1 := errTest.Yeet("")
	if err1.Help() != "" {
		t.Fatal("error help message must be empty")
	}

	err2 := errTestHelp.Yeet("")
	if err2.Help() != "check article 31.40.m" {
		t.Fatal("error help message does not match expected")
	}
}

func TestErrorDefined_Yeet(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet("")

	if err.Code() != "err_test" {
		t.Fatal("err does not have the right code")
	}

	unwrapErr1 := err.Unwrap()
	unwrapErr2 := errors.Unwrap(err)

	if !(unwrapErr1 == nil && unwrapErr2 == nil) {
		t.Fatal("unwrapped errors must be nil")
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
		err := errTest.Wrap(errors.New("failed dial target host"), "")
		if err == nil {
			t.Fatal("err cannot be nil after wrap")
		}

		if err.Code() != "err_test" {
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
		err := errTest.Wrap(parent, "")

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

	t.Run("explain nil *oops.Error", func(t *testing.T) {
		err := Explain(testReturnNilOopsError(), "foo bar baz")
		if err != nil {
			t.Fatal("explain must not create error from nil")
		}
	})

	t.Run("explain new error", func(t *testing.T) {
		err := errors.New("fiz biz")
		errExplained1 := Explain(err, "foo bar")
		errExplained2 := Explain(err, "bar foo")

		if errExplained1 != nil && errExplained1.Code() != "unexpected" {
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

	t.Run("yeet simple", func(t *testing.T) {
		errYeet := errTest.Yeet("")
		j, err := errYeet.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"err_test\",\"type\":\"test\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("yeet help", func(t *testing.T) {
		errYeet := errTestHelp.Yeet("")
		j, err := errYeet.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"err_test_help\",\"type\":\"test\",\"help\":\"check article 31.40.m\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("wrap", func(t *testing.T) {
		errWrap := errTest.Wrap(errors.New("foo bar"), "")

		j, err := errWrap.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"err_test\",\"type\":\"test\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("explain", func(t *testing.T) {
		errExplained := errTest.Yeet("explain 1")

		j, err := errExplained.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"err_test\",\"type\":\"test\",\"explain\":\"explain 1\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("explain x3", func(t *testing.T) {
		errExplained := errTest.Yeet("explain 1")
		_ = Explain(errExplained, "explain 2")
		_ = Explain(errExplained, "explain 3")

		j, err := errExplained.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"err_test\",\"type\":\"test\",\"explain\":\"explain 1, explain 2, explain 3\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})

	t.Run("wrap explain empty", func(t *testing.T) {
		errExplained := errTest.Wrap(errors.New("foo bar"), "explain 1")

		j, err := errExplained.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}

		if string(j) != "{\"code\":\"err_test\",\"type\":\"test\",\"explain\":\"explain 1\"}" {
			t.Fatal("error json does not match expected", string(j))
		}
	})
}

func errorsIs(err, target error) bool {
	return errors.Is(err, target)
}

func TestError_Is(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		err := errTest.Yeet("")

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
		err := errTest.Wrap(parent, "")

		if !errors.Is(err, parent) {
			t.Fatal("error must match parent on Is")
		}

		if !err.Is(parent) {
			t.Fatal("error must match parent on Is")
		}
	})

	t.Run("nil parent", func(t *testing.T) {
		parent := errors.New("foo bar")
		err := errTest.Wrap(nil, "")

		if errors.Is(err, parent) {
			t.Fatal("error must not match parent on Is")
		}

		if err.Is(parent) {
			t.Fatal("error must not match parent on Is")
		}
	})

	t.Run("defined", func(t *testing.T) {
		err := errTest.Yeet("")

		if !errors.Is(err, errTest) {
			t.Fatal("expected err to be defined err")
		}
	})

	t.Run("defined shortcut", func(t *testing.T) {
		err := errTest.Yeet("")

		if !err.Is(errTest) {
			t.Fatal("expected err to be defined err")
		}
	})

	t.Run("defined wrap", func(t *testing.T) {
		var defErr = Define().Code("not_found").Type("not_found")

		err := defErr.Wrap(sql.ErrNoRows, "could not be found in db")

		if !errorsIs(err, sql.ErrNoRows) {
			t.Fatal("expected err to be defined err")
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

func Test_stack(t *testing.T) {
	t.Run("no stack", func(t *testing.T) {
		err := errTest.Yeet("")
		if err.Trace() != nil {
			t.Fatal("no stack error should have no trace")
		}
	})

	t.Run("normal depth", func(t *testing.T) {
		err := errTestTrace.Yeet("")
		if err.Trace() == nil {
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
	return errTestExplainNested.Yeet("source not found")
}

func Test_ExplainFmt(t *testing.T) {
	t.Parallel()

	t.Run("yeet fmt", func(t *testing.T) {
		err := errTest.Yeet("foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})

	t.Run("wrap fmt", func(t *testing.T) {
		err := errTest.Wrap(errors.New("fiz"), "foo %s", "bar")
		if explain := err.Explanation(); explain != "foo bar" {
			t.Fatal("unexpected fmt explain: ", explain)
		}
	})

	t.Run("explain", func(t *testing.T) {
		err := errors.New("new")
		out := Explain(err, "foo %s", "bar")
		msg := out.Error()

		if msg != "[unexpected] unexpected : foo bar" {
			t.Fatalf("unexpected error message('%s')", msg)
		}
	})
}

func Test_returnNil(t *testing.T) {
	t.Parallel()

	if err := testReturnNilOopsError(); err != nil {
		t.Fatal("should not have checked err != nil as true")
	}
}

func TestError_String(t *testing.T) {
	err := errTest.Yeet("foobar")
	_ = Explain(err, "fiz")
	_ = Explain(err, "fuz")

	s := err.String()
	if s != "[test] err_test : foobar, fiz, fuz" {
		t.Fatalf("error string does not match expectations('%s')", s)
	}
}

func BenchmarkError_String(b *testing.B) {
	b.ReportAllocs()

	err := errTest.Yeet("")

	for iter := 0; iter <= b.N; iter++ {
		_ = err.String()
	}
}

func Test_Defer(t *testing.T) {
	t.Parallel()

	t.Run("no fail", func(t *testing.T) {
		err := testDefer("", "")
		if err != nil {
			t.Fatalf("got unexpected error(%s)", err.Error())
		}
	})

	t.Run("fail arg1", func(t *testing.T) {
		err := testDefer("fail", "")
		oopsErr, ok := err.(*Error)
		if !ok {
			t.Fatalf("expected oops.Error, got something else (%+v)", err)
		}

		got := oopsErr.Code()
		wanted := errTest.code

		if got != wanted {
			t.Fatalf("non-matching error codes, got '%s' but wanted '%s'", got, wanted)
		}

		got = oopsErr.Explanation()
		wanted = "failed test defer do 1, failed test defer with args1='fail' and arg2=''"

		if got != wanted {
			t.Fatalf("non-matching error explanations, got '%s' but wanted '%s'", got, wanted)
		}
	})

	t.Run("fail arg2", func(t *testing.T) {
		err := testDefer("", "fail")
		oopsErr, ok := err.(*Error)
		if !ok {
			t.Fatalf("expected oops.Error, got something else (%+v) (%T)", err, err)
		}

		got := oopsErr.Code()
		wanted := ErrUnexpected.code

		if got != wanted {
			t.Fatalf("non-matching error codes, got '%s' but wanted '%s'", got, wanted)
		}

		got = oopsErr.Explanation()
		wanted = "deferred error, failed test defer with args1='' and arg2='fail'"

		if got != wanted {
			t.Fatalf("non-matching error explanations, got '%s' but wanted '%s'", got, wanted)
		}

		parent := oopsErr.Unwrap()
		if parent.Error() != "failed test defer do 2" {
			t.Fatalf("bad parent message('%s')", parent.Error())
		}
	})

	t.Run("nil", func(t *testing.T) {
		Defer(nil, "foo %s", "bar")
	})
}

func testDefer(arg1, arg2 string) (err error) {
	defer Defer(&err, "failed test defer with args1='%s' and arg2='%s'", arg1, arg2)

	err = testDeferDo1(arg1)
	if err != nil {
		return err
	}

	err = testDeferDo2(arg2)
	if err != nil {
		return err
	}

	return nil
}

func testDeferDo1(arg string) error {
	if arg == "fail" {
		return errTest.Yeet("failed test defer do 1")
	}

	return nil
}

func testDeferDo2(arg string) error {
	if arg == "fail" {
		return errors.New("failed test defer do 2")
	}

	return nil
}

func TestError_StatusCode(t *testing.T) {
	t.Parallel()

	t.Run("no status code", func(t *testing.T) {
		err := errTest.Yeet("")
		if err.StatusCode() != 0 {
			t.Fatalf("expected status code 0, got %d", err.StatusCode())
		}
	})

	t.Run("status code 418", func(t *testing.T) {
		err := errTestStatusCode.Yeet("")
		if err.StatusCode() != 418 {
			t.Fatalf("expected status code 418, got %d", err.StatusCode())
		}
	})
}

func Test_errorGroup(t *testing.T) {
	t.Parallel()

	var (
		group = Define().Code("ignored").Type("test_group").StatusCode(500).Group()

		err1 = group.Code("1")
		err2 = group.Code("2").Type("overwrite")
	)

	err1spawn := err1.Yeet("")
	err2spawn := err2.Yeet("")

	if err1spawn.StatusCode() != 500 {
		t.Errorf("expected status code 500, got %d", err1spawn.StatusCode())
	}
	if err2spawn.StatusCode() != 500 {
		t.Errorf("expected status code 500, got %d", err2spawn.StatusCode())
	}

	if err1spawn.Code() != "1" {
		t.Errorf("expected code 1, got %s", err1spawn.Code())
	}
	if err2spawn.Code() != "2" {
		t.Errorf("expected code 2, got %s", err2spawn.Code())
	}

	if err1spawn.Type() != "test_group" {
		t.Errorf("expected type test_group, got %s", err1spawn.Type())
	}
	if err2spawn.Type() != "overwrite" {
		t.Errorf("expected type overwrite, got %s", err2spawn.Type())
	}
}

func Test_errorGroup_PrefixCode(t *testing.T) {
	t.Parallel()

	group := Define().Type("prefix").Group().PrefixCode("yes_")
	err := group.Code("foo").Yeet("")

	if err.Code() != "yes_foo" {
		t.Errorf("expected code 'yes_foo', got %s", err.Code())
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
		return Explain(err, "nested error 1")
	}

	return nil
}

func benchmarkNested2(original error) error {
	if err := benchmarkNested3(original); err != nil {
		return Explain(err, "nested error 2")
	}

	return nil
}

func benchmarkNested3(original error) error {
	err := benchmarkNested4(original)
	if err != nil {
		return Explain(err, "nested error 3")
	}

	return nil
}

func benchmarkNested4(original error) error {
	return errTestBenchmark.Wrap(original, "benchmarkNested4 returned wrapped original error")
}

func TestError_Explain(t *testing.T) {
	t.Parallel()

	err := errTest.Yeet("foo bar").Explain("baz")

	if err.Explanation() != "foo bar, baz" {
		t.Errorf("expected explanation 'foo bar, baz', got %s", err.Explanation())
	}

	_ = err.Explain("id=%d", 123)

	if err.Explanation() != "foo bar, baz, id=123" {
		t.Errorf("expected explanation 'foo bar, baz, id=123', got %s", err.Explanation())
	}
}

func TestError_Err(t *testing.T) {
	t.Parallel()

	t.Run("the problem", func(t *testing.T) {
		var err error
		err = Explain(nil, "foobar")
		if err == nil {
			t.Errorf("expected nil to not match nil")
		}
	})

	t.Run("the solution", func(t *testing.T) {
		var err error
		err = Explain(nil, "foobar").Err()
		if err != nil {
			t.Errorf("expected nil to match nil")
		}
	})
}
