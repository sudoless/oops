package stack

import (
	"strings"
	"testing"
)

//go:noinline
func recurse(depth int, fn func()) {
	if depth == 0 {
		fn()
		return
	}
	recurse(depth-1, fn)
}

func TestStack(t *testing.T) {
	t.Parallel()

	t.Run("skip 1 starts at the caller", func(t *testing.T) {
		t.Parallel()
		frames := Stack(1)
		if len(frames) == 0 {
			t.Fatal("no frames")
		}
		if !strings.HasSuffix(frames[0], ": TestStack.func1") {
			t.Fatalf("frames[0] = %q", frames[0])
		}
		if !strings.Contains(frames[0], "stack_test.go:") {
			t.Fatalf("frames[0] missing file:line: %q", frames[0])
		}
	})

	t.Run("skip 0 starts at Stack", func(t *testing.T) {
		t.Parallel()
		frames := Stack(0)
		if len(frames) == 0 || !strings.HasSuffix(frames[0], ": Stack") {
			t.Fatalf("frames = %q", frames)
		}
	})

	t.Run("capped at 32 frames", func(t *testing.T) {
		t.Parallel()
		var frames []string
		recurse(100, func() { frames = Stack(1) })
		if len(frames) != 32 {
			t.Fatalf("got %d frames, want 32", len(frames))
		}
		if !strings.HasSuffix(frames[1], ": recurse") {
			t.Fatalf("frames[1] = %q", frames[1])
		}
	})
}

func TestFunction(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct{ in, want string }{
		{"", "???"},
		{"main.main", "main"},
		{"go.sdls.io/oops/v2/pkg/oops.(*ErrorDefinition).Yeet", "(*ErrorDefinition).Yeet"},
		{"go.sdls.io/oops/v2/pkg/oops_test.TestTrace.func1", "TestTrace.func1"},
		{"example.com/pkg.T·method", "T.method"},
	} {
		if got := function(tc.in); got != tc.want {
			t.Errorf("function(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
