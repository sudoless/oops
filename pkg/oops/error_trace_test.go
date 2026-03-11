package oops_test

import (
	"strings"
	"testing"

	"go.sdls.io/oops/pkg/oops"
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
		if !strings.Contains(first, "(0x") {
			t.Fatalf("frame missing hex PC: %q", first)
		}
		if !strings.Contains(first, "): ") {
			t.Fatalf("frame missing ): separator: %q", first)
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
