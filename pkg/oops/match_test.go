package oops_test

import (
	"testing"

	"go.sdls.io/oops/pkg/oops"
)

func TestMatchers(t *testing.T) {
	t.Parallel()

	t.Run("ByCause", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Causes(oops.CauseNotFound).Yeet()
		if !oops.Match(err, oops.ByCause(oops.CauseNotFound)) {
			t.Fatal("expected match")
		}
		if oops.Match(err, oops.ByCause(oops.CauseAuth)) {
			t.Fatal("expected no match")
		}
	})

	t.Run("ByAction", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Actions(oops.ActionRetry).Yeet()
		if !oops.Match(err, oops.ByAction(oops.ActionRetry)) {
			t.Fatal("expected match")
		}
	})

	t.Run("ByCode", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("auth.expired").Yeet()
		if !oops.Match(err, oops.ByCode("auth.expired")) {
			t.Fatal("expected match")
		}
		if oops.Match(err, oops.ByCode("other")) {
			t.Fatal("expected no match")
		}
	})

	t.Run("ByDefinition", func(t *testing.T) {
		t.Parallel()
		def := oops.Define("test")
		err := def.Yeet()
		if !oops.Match(err, oops.ByDefinition(def)) {
			t.Fatal("expected match")
		}
	})

	t.Run("ByDefinition with inherits", func(t *testing.T) {
		t.Parallel()
		base := oops.Define("base")
		child := oops.Define("child").Inherits(base)
		err := child.Yeet()
		if !oops.Match(err, oops.ByDefinition(base)) {
			t.Fatal("expected match via inheritance")
		}
	})

	t.Run("All", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Causes(oops.CauseAuth).Actions(oops.ActionRetry).Yeet()
		m := oops.All(oops.ByCause(oops.CauseAuth), oops.ByAction(oops.ActionRetry))
		if !oops.Match(err, m) {
			t.Fatal("expected All to match")
		}

		m2 := oops.All(oops.ByCause(oops.CauseAuth), oops.ByAction(oops.ActionAbort))
		if oops.Match(err, m2) {
			t.Fatal("expected All to fail with missing action")
		}
	})

	t.Run("Any", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Causes(oops.CauseAuth).Yeet()
		m := oops.Any(oops.ByCause(oops.CauseNotFound), oops.ByCause(oops.CauseAuth))
		if !oops.Match(err, m) {
			t.Fatal("expected Any to match")
		}

		m2 := oops.Any(oops.ByCause(oops.CauseNotFound), oops.ByCause(oops.CauseTimeout))
		if oops.Match(err, m2) {
			t.Fatal("expected Any to fail")
		}
	})

	t.Run("Not", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Causes(oops.CauseAuth).Yeet()
		if oops.Match(err, oops.Not(oops.ByCause(oops.CauseAuth))) {
			t.Fatal("expected Not to invert")
		}
		if !oops.Match(err, oops.Not(oops.ByCause(oops.CauseNotFound))) {
			t.Fatal("expected Not to invert")
		}
	})

	t.Run("Match nil", func(t *testing.T) {
		t.Parallel()
		if oops.Match(nil, oops.ByCause(oops.CauseAuth)) {
			t.Fatal("expected false for nil")
		}
	})
}

func TestPrebuiltMatchers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		matcher oops.Matcher
		cause   string
		action  string
	}{
		{"IsRetryable", oops.IsRetryable, "", oops.ActionRetry},
		{"IsAbort", oops.IsAbort, "", oops.ActionAbort},
		{"IsFatal", oops.IsFatal, "", oops.ActionFatal},
		{"IsAuth", oops.IsAuth, oops.CauseAuth, ""},
		{"IsNotFound", oops.IsNotFound, oops.CauseNotFound, ""},
		{"IsTimeout", oops.IsTimeout, oops.CauseTimeout, ""},
		{"IsInternal", oops.IsInternal, oops.CauseInternal, ""},
		{"IsRateLimit", oops.IsRateLimit, oops.CauseRateLimit, ""},
		{"IsConflict", oops.IsConflict, oops.CauseConflict, ""},
		{"IsValidation", oops.IsValidation, oops.CauseValidation, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			def := oops.Define("test")
			if tt.cause != "" {
				def = def.Causes(tt.cause)
			}
			if tt.action != "" {
				def = def.Actions(tt.action)
			}
			err := def.Yeet()
			if !oops.Match(err, tt.matcher) {
				t.Fatalf("expected %s to match", tt.name)
			}
		})
	}
}
