package oops_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

var (
	errWalkA = oops.Define("walk_a")
	errWalkB = oops.Define("walk_b")
)

// multiUnwrapError is a foreign error with Unwrap() []error that keeps nil children.
type multiUnwrapError struct{ errs []error }

func (*multiUnwrapError) Error() string     { return "multi" }
func (m *multiUnwrapError) Unwrap() []error { return m.errs }

// singleUnwrapError is a foreign error with Unwrap() error.
type singleUnwrapError struct{ err error }

func (singleUnwrapError) Error() string   { return "single" }
func (s singleUnwrapError) Unwrap() error { return s.err }

// selfUnwrapError is a foreign error that unwraps to itself.
type selfUnwrapError struct{}

func (*selfUnwrapError) Error() string   { return "self" }
func (s *selfUnwrapError) Unwrap() error { return s }

type walkNode struct {
	depth int
	err   error
}

func walkAll(err error) []walkNode {
	var out []walkNode
	for depth, e := range oops.Walk(err) {
		out = append(out, walkNode{depth, e})
	}
	return out
}

func assertWalk(t *testing.T, got, want []walkNode) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("walked %d nodes, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i].depth != want[i].depth || got[i].err != want[i].err { //nolint:errorlint // identity check
			t.Fatalf("node %d = (%d, %v), want (%d, %v)", i, got[i].depth, got[i].err, want[i].depth, want[i].err)
		}
	}
}

func TestWalk(t *testing.T) {
	t.Parallel()

	t.Run("single", func(t *testing.T) {
		t.Parallel()
		err := errWalkA.Yeet()
		assertWalk(t, walkAll(err), []walkNode{{0, err}})
	})

	t.Run("chain", func(t *testing.T) {
		t.Parallel()
		inner := errWalkB.Wrap(io.EOF)
		outer := errWalkA.Wrap(inner)
		assertWalk(t, walkAll(outer), []walkNode{{0, outer}, {1, inner}, {2, io.EOF}})
	})

	t.Run("fan in order", func(t *testing.T) {
		t.Parallel()
		c1, c2, c3 := errWalkB.Yeet(), errWalkB.Yeet(), errWalkB.Yeet()
		root := errWalkA.Yeet().Nest(c1, c2, c3)
		assertWalk(t, walkAll(root), []walkNode{{0, root}, {1, c1}, {1, c2}, {1, c3}})
	})

	t.Run("nested pre-order", func(t *testing.T) {
		t.Parallel()
		a1, a2, b1 := errWalkB.Yeet(), errWalkB.Yeet(), errWalkB.Yeet()
		a := errWalkA.Yeet().Nest(a1, a2)
		b := errWalkA.Yeet().Nest(b1)
		root := errWalkA.Yeet().Nest(a, b)
		assertWalk(t, walkAll(root), []walkNode{
			{0, root}, {1, a}, {2, a1}, {2, a2}, {1, b}, {2, b1},
		})
	})

	t.Run("foreign nodes are yielded and followed", func(t *testing.T) {
		t.Parallel()
		dial := fmt.Errorf("dial: %w", io.ErrUnexpectedEOF)
		inner := errors.Join(io.EOF, io.ErrClosedPipe)
		query := fmt.Errorf("q: %w", inner)
		b := errWalkB.Wrap(query)
		plain := errors.New("plain")
		join := errors.Join(dial, b, plain)
		root := errWalkA.Wrap(join)
		assertWalk(t, walkAll(root), []walkNode{
			{0, root},
			{1, join},
			{2, dial},
			{3, io.ErrUnexpectedEOF},
			{2, b},
			{3, query},
			{4, inner},
			{5, io.EOF},
			{5, io.ErrClosedPipe},
			{2, plain},
		})
	})

	t.Run("nil and typed-nil roots yield nothing", func(t *testing.T) {
		t.Parallel()
		if got := walkAll(nil); len(got) != 0 {
			t.Fatalf("nil: %v", got)
		}
		if got := walkAll((*oops.Error)(nil)); len(got) != 0 {
			t.Fatalf("typed nil: %v", got)
		}
	})

	t.Run("nil and typed-nil children are skipped", func(t *testing.T) {
		t.Parallel()
		leaf := errWalkB.Yeet()
		single := singleUnwrapError{(*oops.Error)(nil)}
		multi := &multiUnwrapError{[]error{nil, (*oops.Error)(nil), single, leaf}}
		assertWalk(t, walkAll(multi), []walkNode{{0, multi}, {1, single}, {1, leaf}})
	})

	t.Run("self-nested error walks one node", func(t *testing.T) {
		t.Parallel()
		err := errWalkA.Yeet()
		_ = err.Nest(err)
		assertWalk(t, walkAll(err), []walkNode{{0, err}})
	})

	t.Run("cycle below the root is cut at the repeated ancestor", func(t *testing.T) {
		t.Parallel()
		a, b := errWalkA.Yeet(), errWalkB.Yeet()
		_ = a.Nest(b)
		_ = b.Nest(a)
		assertWalk(t, walkAll(a), []walkNode{{0, a}, {1, b}})
	})

	t.Run("cycle deeper than 32 ancestors is cut", func(t *testing.T) {
		t.Parallel()
		chain := make([]*oops.Error, 40)
		for i := range chain {
			chain[i] = errWalkA.Yeet()
		}
		want := make([]walkNode, 0, len(chain))
		for i, e := range chain {
			_ = e.Nest(chain[(i+1)%len(chain)])
			want = append(want, walkNode{i, e})
		}
		assertWalk(t, walkAll(chain[0]), want)
	})

	t.Run("shared child is not a cycle", func(t *testing.T) {
		t.Parallel()
		child := errWalkB.Yeet()
		root := errWalkA.Yeet().Nest(child, child)
		assertWalk(t, walkAll(root), []walkNode{{0, root}, {1, child}, {1, child}})
	})

	t.Run("stops after 1024 nodes", func(t *testing.T) {
		t.Parallel()
		children := make([]error, 2000)
		for i := range children {
			children[i] = errWalkB.Yeet()
		}
		root := errWalkA.Yeet().Nest(children...)
		if got := len(walkAll(root)); got != 1024 {
			t.Fatalf("fan: walked %d nodes, want 1024", got)
		}
		if got := len(walkAll(&selfUnwrapError{})); got != 1024 {
			t.Fatalf("foreign self-cycle: walked %d nodes, want 1024", got)
		}
	})

	t.Run("early break stops the walk", func(t *testing.T) {
		t.Parallel()
		c1, c2 := errWalkB.Yeet(), errWalkB.Yeet()
		root := errWalkA.Yeet().Nest(c1.Nest(io.EOF), c2)
		var got []walkNode
		for depth, e := range oops.Walk(root) {
			got = append(got, walkNode{depth, e})
			if e == c1 { //nolint:errorlint // identity check
				break
			}
		}
		assertWalk(t, got, []walkNode{{0, root}, {1, c1}})
	})

	t.Run("early break below depth 1 stops the walk", func(t *testing.T) {
		t.Parallel()
		c1, c2 := errWalkB.Yeet(), errWalkB.Yeet()
		root := errWalkA.Yeet().Nest(c1.Nest(io.EOF), c2)
		var got []walkNode
		for depth, e := range oops.Walk(root) {
			got = append(got, walkNode{depth, e})
			if e == io.EOF { //nolint:errorlint // identity check
				break
			}
		}
		assertWalk(t, got, []walkNode{{0, root}, {1, c1}, {2, io.EOF}})
	})
}

func TestWalk_Allocations(t *testing.T) { //nolint:paralleltest // AllocsPerRun panics in a parallel test
	if raceEnabled {
		t.Skip("the race detector changes allocation counts")
	}

	for _, tree := range benchTrees() {
		allocs := testing.AllocsPerRun(100, func() {
			n := 0
			for depth := range oops.Walk(tree.err) {
				n += depth
			}
			benchSinkInt = n
		})
		if allocs != 0 {
			t.Errorf("%s: %v allocations per full walk, want 0", tree.name, allocs)
		}
	}
}
