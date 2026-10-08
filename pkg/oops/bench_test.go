package oops_test

import (
	"errors"
	"fmt"
	"io"
	"testing"

	"go.sdls.io/oops/v2/pkg/oops"
)

var (
	benchDef       = oops.Define("bench").Causes(oops.CauseIO).Message("bench failed")
	benchDefTraced = oops.Define("bench.traced").Traced()

	benchSinkErr *oops.Error
	benchSinkStr string
	benchSinkInt int
)

func BenchmarkYeet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchSinkErr = benchDef.Yeet()
	}
}

func BenchmarkWrap(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		benchSinkErr = benchDef.Wrap(io.EOF)
	}
}

// atStackDepth calls fn with depth extra frames on the stack.
//
//go:noinline
func atStackDepth(depth int, fn func()) {
	if depth == 0 {
		fn()
		return
	}
	atStackDepth(depth-1, fn)
}

// BenchmarkYeetTraced measures trace capture with the caller at increasing
// stack depths: the cost must not grow with the depth.
func BenchmarkYeetTraced(b *testing.B) {
	for _, bc := range []struct {
		name  string
		depth int
	}{
		{"depth_10", 10},
		{"depth_100", 100},
		{"depth_1000", 1000},
	} {
		b.Run(bc.name, func(b *testing.B) {
			atStackDepth(bc.depth, func() {
				b.ReportAllocs()
				for b.Loop() {
					benchSinkErr = benchDefTraced.Yeet()
				}
			})
		})
	}
}

func BenchmarkErrorString(b *testing.B) {
	for _, bc := range []struct {
		name string
		err  *oops.Error
	}{
		{"code", oops.Define("bench").Yeet()},
		{"message_explanation", benchDef.Yeetf("item %d", 7)},
		{"wrapped", benchDef.Wrapf(errors.New("dial tcp: refused"), "connect")},
	} {
		b.Run(bc.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				benchSinkStr = bc.err.Error()
			}
		})
	}
}

// benchTrees returns error trees of the shapes Walk is measured on.
func benchTrees() []struct {
	name string
	err  error
} {
	defA := oops.Define("bench_a").Causes(oops.CauseIO)
	defB := oops.Define("bench_b").Causes(oops.CauseNotFound)
	defV := oops.Define("bench_validation").Causes(oops.CauseValidation)

	single := defA.Yeet()

	chain := error(io.EOF)
	for i := range 5 {
		chain = defA.Wrapf(chain, "layer %d", i)
	}

	finish, add := defV.Collect()
	for i := range 50 {
		add(defB.Yeetf("field %d", i), "f%d", i)
	}
	fan := finish()

	outerFinish, outerAdd := defV.Collect()
	for i := range 10 {
		innerFinish, innerAdd := defV.Collect()
		for j := range 10 {
			innerAdd(defB.Yeet(), "f%d", j)
		}
		outerAdd(innerFinish(), "item%d", i)
	}
	nested := outerFinish()

	mixed := defA.Wrap(errors.Join(
		fmt.Errorf("dial: %w", io.ErrUnexpectedEOF),
		defB.Wrap(fmt.Errorf("q: %w", errors.Join(io.EOF, io.ErrClosedPipe))),
		errors.New("plain"),
	))

	return []struct {
		name string
		err  error
	}{
		{"single_1", single},
		{"chain_6", chain},
		{"fan_51", fan},
		{"nested_111", nested},
		{"mixed_10", mixed},
	}
}

func BenchmarkWalk(b *testing.B) {
	for _, tree := range benchTrees() {
		b.Run(tree.name, func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				n := 0
				for depth, err := range oops.Walk(tree.err) {
					if _, ok := oops.Native(err); ok {
						n += depth
					}
					n++
				}
				benchSinkInt = n
			}
		})
		// stop at the first node below the root
		b.Run(tree.name+"/break", func(b *testing.B) {
			b.ReportAllocs()
			for b.Loop() {
				n := 0
				for depth := range oops.Walk(tree.err) {
					n++
					if depth >= 1 {
						break
					}
				}
				benchSinkInt = n
			}
		})
	}
}
