package oops_test

import (
	"database/sql"
	"errors"
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestError_Is(t *testing.T) {
	t.Parallel()

	errorsIs := func(err, target error) bool {
		return errors.Is(err, target)
	}

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
		err := oops.NilErr

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

	t.Run("defined", func(t *testing.T) {
		err := errTest.Yeet()

		if !errors.Is(err, errTest) {
			t.Fatal("expected err to be defined err")
		}
	})

	t.Run("defined shortcut", func(t *testing.T) {
		err := errTest.Yeet()

		if !err.Is(errTest) {
			t.Fatal("expected err to be defined err")
		}
	})

	t.Run("defined wrap", func(t *testing.T) {
		defErr := oops.Define("code", "test.not_found")

		err := defErr.Wrapf(sql.ErrNoRows, "could not be found in db")

		if !errorsIs(err, sql.ErrNoRows) {
			t.Fatal("expected err to be defined err")
		}
	})
}

func TestError_As(t *testing.T) {
	t.Parallel()

	t.Run("basic", func(t *testing.T) {
		someErr := error(errTest.Wrapf(errors.New("fiz"), "foobar"))

		var oopsErr oops.Error
		if !errors.As(someErr, &oopsErr) {
			t.Fatal("expected to find oops error")
		}

		if oopsErr == nil {
			t.Fatal("expected oops error to be non-nil")
		}

		if oopsErr.Explanation() != "foobar" {
			t.Fatal("expected oops error to have explanation")
		}
	})
}

func TestError_error(t *testing.T) {
	t.Parallel()

	t.Run("Explainf", func(t *testing.T) {
		var err error //nolint:gosimple
		err = oops.Explainf(nil, "foobar")
		if err != nil {
			t.Errorf("expected nil to match nil")
		}
	})

	t.Run("return nil", func(t *testing.T) {
		returnNil := func() oops.Error {
			return nil
		}

		if err := returnNil(); err != nil {
			t.Fatal("should not have checked err != nil as true")
		}
	})
}
