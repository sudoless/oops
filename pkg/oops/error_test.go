package oops_test

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

func TestError_Accessors(t *testing.T) {
	t.Parallel()

	t.Run("Definition", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if err.Definition() != def {
			t.Fatal("expected definition to match")
		}
	})

	t.Run("Code", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("my.code").Yeet()
		if err.Code() != "my.code" {
			t.Fatalf("got %q", err.Code())
		}
	})

	t.Run("Message", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Message("hello").Yeet()
		if err.Message() != "hello" {
			t.Fatalf("got %q", err.Message())
		}
	})
}

// errorMethods calls every exported *Error method once.
var errorMethods = map[string]func(*oops.Error){
	"Definition":  func(err *oops.Error) { _ = err.Definition() },
	"Code":        func(err *oops.Error) { _ = err.Code() },
	"Message":     func(err *oops.Error) { _ = err.Message() },
	"Explanation": func(err *oops.Error) { _ = err.Explanation() },
	"Causes":      func(err *oops.Error) { _ = err.Causes() },
	"Actions":     func(err *oops.Error) { _ = err.Actions() },
	"Fields":      func(err *oops.Error) { _ = err.Fields() },
	"Get":         func(err *oops.Error) { _, _ = err.Get("key") },
	"HasCause":    func(err *oops.Error) { _ = err.HasCause(oops.CauseIO) },
	"HasAction":   func(err *oops.Error) { _ = err.HasAction(oops.ActionRetry) },
	"Error":       func(err *oops.Error) { _ = err.Error() },
	"Unwrap":      func(err *oops.Error) { _ = err.Unwrap() },
	"Is":          func(err *oops.Error) { _ = err.Is(oops.ErrTODO) },
	"Explainf":    func(err *oops.Error) { _ = err.Explainf("x %d", 1) },
	"Set":         func(err *oops.Error) { _ = err.Set("key", 1) },
	"AddCauses":   func(err *oops.Error) { _ = err.AddCauses(oops.CauseIO) },
	"SetActions":  func(err *oops.Error) { _ = err.SetActions(oops.ActionRetry) },
	"Nest":        func(err *oops.Error) { _ = err.Nest(errors.New("x")) },
	"Pathf":       func(err *oops.Error) { _ = err.Pathf("p/%d", 1) },
	"Path":        func(err *oops.Error) { _ = err.Path() },
	"PathArgs":    func(err *oops.Error) { _ = err.PathArgs() },
	"Trace":       func(err *oops.Error) { _ = err.Trace() },
	"Clone":       func(err *oops.Error) { _ = err.Clone() },
}

func TestError_NilAndZeroValue(t *testing.T) {
	t.Parallel()

	t.Run("table covers every exported method", func(t *testing.T) {
		t.Parallel()
		typ := reflect.TypeFor[*oops.Error]()
		for idx := range typ.NumMethod() {
			name := typ.Method(idx).Name
			if _, ok := errorMethods[name]; !ok {
				t.Errorf("method %s missing from errorMethods", name)
			}
		}
	})

	receivers := map[string]func() *oops.Error{
		"nil":        func() *oops.Error { return nil },
		"zero value": func() *oops.Error { return &oops.Error{} },
	}
	for recvName, newRecv := range receivers {
		for methodName, call := range errorMethods {
			t.Run(recvName+"/"+methodName, func(t *testing.T) {
				t.Parallel()
				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("panicked: %v", r)
					}
				}()
				call(newRecv())
			})
		}
	}

	t.Run("zero value reads as empty", func(t *testing.T) {
		t.Parallel()
		err := &oops.Error{}
		if err.Code() != "" || err.Message() != "" || err.Definition() != nil {
			t.Fatalf("got code %q, message %q, definition %v", err.Code(), err.Message(), err.Definition())
		}
		if err.Is(oops.ErrTODO) {
			t.Fatal("zero value must not match a definition")
		}
		if _, ok := oops.As(err, oops.ErrTODO); ok {
			t.Fatal("zero value must not match in As")
		}
	})

	t.Run("nil mutators return nil", func(t *testing.T) {
		t.Parallel()
		var err *oops.Error
		if err.Explainf("x") != nil || err.Set("k", 1) != nil || err.AddCauses("c") != nil ||
			err.SetActions("a") != nil || err.Nest(errors.New("x")) != nil || err.Pathf("p") != nil ||
			err.Clone() != nil {
			t.Fatal("expected nil from every mutator on a nil receiver")
		}
	})
}

func TestError_Clone(t *testing.T) {
	t.Parallel()

	def := oops.Define("test").Causes(oops.CauseIO).Actions(oops.ActionRetry)
	child := oops.Define("child").Yeet()
	orig := def.Yeetf("first").Set("key", 1).Pathf("item/%d", 1).Nest(child)

	clone := orig.Clone()
	if clone == orig {
		t.Fatal("Clone returned the receiver")
	}
	if clone.Definition() != def || !errors.Is(clone, def) {
		t.Fatal("clone lost its definition")
	}
	if clone.Unwrap()[0] != child {
		t.Fatal("clone must share wrapped children")
	}

	_ = clone.Explainf("second").Set("key", 2).AddCauses(oops.CauseTimeout).
		SetActions(oops.ActionAbort).Pathf("other").Nest(errors.New("clone child"))
	_ = orig.AddCauses(oops.CauseAuth).Nest(errors.New("orig child"))

	if orig.Explanation() != "first" {
		t.Errorf("orig explanation = %q", orig.Explanation())
	}
	if v, _ := orig.Get("key"); v != 1 {
		t.Errorf("orig field = %v", v)
	}
	if got := orig.Causes(); !slices.Equal(got, []string{oops.CauseIO, oops.CauseAuth}) {
		t.Errorf("orig causes = %v", got)
	}
	if got := orig.Actions(); !slices.Equal(got, []string{oops.ActionRetry}) {
		t.Errorf("orig actions = %v", got)
	}
	if orig.Path() != "item/1" || !slices.Equal(orig.PathArgs(), []any{1}) {
		t.Errorf("orig path = %q %v", orig.Path(), orig.PathArgs())
	}
	if got := orig.Unwrap(); len(got) != 2 || got[1].Error() != "orig child" {
		t.Errorf("orig wrapped = %v", got)
	}

	if clone.Explanation() != "first, second" {
		t.Errorf("clone explanation = %q", clone.Explanation())
	}
	if v, _ := clone.Get("key"); v != 2 {
		t.Errorf("clone field = %v", v)
	}
	if got := clone.Causes(); !slices.Equal(got, []string{oops.CauseIO, oops.CauseTimeout}) {
		t.Errorf("clone causes = %v", got)
	}
	if got := clone.Actions(); !slices.Equal(got, []string{oops.ActionAbort}) {
		t.Errorf("clone actions = %v", got)
	}
	if clone.Path() != "other" || clone.PathArgs() != nil {
		t.Errorf("clone path = %q %v", clone.Path(), clone.PathArgs())
	}
	if got := clone.Unwrap(); len(got) != 2 || got[1].Error() != "clone child" {
		t.Errorf("clone wrapped = %v", got)
	}

	t.Run("no shared backing arrays", func(t *testing.T) {
		t.Parallel()
		src := def.Yeet().Pathf("item/%d", 1).Nest(child)
		dup := src.Clone()
		dup.Causes()[0] = "mutated"
		dup.Actions()[0] = "mutated"
		dup.Unwrap()[0] = errors.New("mutated")
		dup.PathArgs()[0] = "mutated"
		if src.Causes()[0] != oops.CauseIO || src.Actions()[0] != oops.ActionRetry ||
			src.Unwrap()[0] != child || src.PathArgs()[0] != 1 {
			t.Fatalf("clone shares storage with its source: %v %v %v %v",
				src.Causes(), src.Actions(), src.Unwrap(), src.PathArgs())
		}
	})
}
