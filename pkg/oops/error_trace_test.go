package oops_test

import (
	"io"
	"regexp"
	"strings"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

var errTraced = oops.Define("trace.test").Traced()

func traceLevel3() *oops.Error { return errTraced.Yeet() }
func traceLevel2() *oops.Error { return traceLevel3() }
func traceLevel1() *oops.Error { return traceLevel2() }

func TestError_Trace(t *testing.T) {
	t.Parallel()

	t.Run("nil when not traced", func(t *testing.T) {
		t.Parallel()
		err := oops.Define("test").Yeet()
		if err.Trace() != nil {
			t.Fatal("expected nil trace for non-traced definition")
		}
	})

	t.Run("non-nil when traced", func(t *testing.T) {
		t.Parallel()
		err := errTraced.Yeet()
		if err.Trace() == nil {
			t.Fatal("expected non-nil trace for traced definition")
		}
	})

	t.Run("frame format", func(t *testing.T) {
		t.Parallel()
		err := errTraced.Yeet()
		frames := err.Trace()
		if len(frames) == 0 {
			t.Fatal("expected at least one frame")
		}
		first := frames[0]
		if !strings.Contains(first, ".go:") {
			t.Fatalf("frame missing .go: file reference: %q", first)
		}
	})

	t.Run("call stack order with 3 stubs", func(t *testing.T) {
		t.Parallel()
		err := traceLevel1()
		frames := err.Trace()
		if len(frames) < 3 {
			t.Fatalf("expected at least 3 frames, got %d: %v", len(frames), frames)
		}

		expected := []string{"traceLevel3", "traceLevel2", "traceLevel1"}
		for idx, want := range expected {
			if !strings.Contains(frames[idx], want) {
				t.Errorf("frames[%d] = %q, want it to contain %q", idx, frames[idx], want)
			}
		}
	})
}

// TestError_TraceEntryPoints checks that the first frame of every trace is the
// code that called oops, whichever entry point created the error.
func TestError_TraceEntryPoints(t *testing.T) {
	t.Parallel()

	native := func(err error) *oops.Error {
		e, ok := oops.Native(err)
		if !ok {
			t.Fatalf("expected an *oops.Error, got %T", err)
		}
		return e
	}

	cases := []struct {
		name string
		call func() *oops.Error
	}{
		{"Yeet", func() *oops.Error { return errTraced.Yeet() }},
		{"Yeetf", func() *oops.Error { return errTraced.Yeetf("x %d", 1) }},
		{"Wrap", func() *oops.Error { return errTraced.Wrap(io.EOF) }},
		{"Wrapf", func() *oops.Error { return errTraced.Wrapf(io.EOF, "x %d", 1) }},
		{"Collect finish", func() *oops.Error {
			finish, add := errTraced.Collect()
			add(oops.Define("child").Yeet(), "")
			return native(finish())
		}},
		{"Collect add foreign", func() *oops.Error {
			finish, add := oops.Define("untraced").Collect()
			add(io.EOF, "")
			return native(native(finish()).Unwrap()[0])
		}},
		{"Nest", func() *oops.Error { return native(oops.Nest(errTraced, io.EOF)) }},
		{"Foreign", func() *oops.Error { return oops.Foreign(io.EOF) }},
		{"Explainf", func() *oops.Error { return native(oops.Explainf(io.EOF, "x")) }},
		{"AddCauses", func() *oops.Error { return native(oops.AddCauses(io.EOF, oops.CauseIO)) }},
		{"Pathf", func() *oops.Error { return native(oops.Pathf(io.EOF, "p")) }},
	}

	caller := regexp.MustCompile(`error_trace_test\.go:\d+ \(0x[0-9a-f]+\): TestError_TraceEntryPoints\.func\d+$`)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			frames := tc.call().Trace()
			if len(frames) == 0 {
				t.Fatal("no trace")
			}
			if !caller.MatchString(frames[0]) {
				t.Fatalf("frames[0] = %q, want the calling test closure", frames[0])
			}
		})
	}
}
