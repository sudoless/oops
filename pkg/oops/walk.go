package oops

import "iter"

// walkLimit caps the number of nodes Walk yields.
const walkLimit = 1024

// Walk returns an iterator over err and every error below it, in pre-order,
// with each node's depth (0 for err). It follows both Unwrap() error and
// Unwrap() []error, so foreign errors are yielded and traversed too. Nil
// nodes, including a typed-nil *Error, are skipped.
//
// An *Error that appears below itself is a cycle: the repeated node and its
// subtree are skipped. The same error reached through two different parents is
// yielded twice. Walk stops after 1024 nodes, which also bounds cycles made of
// foreign errors.
func Walk(err error) iter.Seq2[int, error] {
	return func(yield func(int, error) bool) {
		n := 0
		var ancestorsBuf [32]*Error
		ancestors := ancestorsBuf[:0]

		var walk func(err error, depth int) bool
		walk = func(err error, depth int) bool {
			if isNil(err) {
				return true
			}

			e, native := err.(*Error) //nolint:errorlint // traversal inspects each node directly
			if native {
				for _, ancestor := range ancestors {
					if ancestor == e {
						return true
					}
				}
			}

			if n++; n > walkLimit {
				return false
			}
			if !yield(depth, err) {
				return false
			}

			if native {
				ancestors = append(ancestors, e)
			}

			more := true
			switch u := err.(type) { //nolint:errorlint // traversal dispatches on the unwrap interfaces
			case interface{ Unwrap() []error }:
				for _, child := range u.Unwrap() {
					if more = walk(child, depth+1); !more {
						break
					}
				}
			case interface{ Unwrap() error }:
				more = walk(u.Unwrap(), depth+1)
			}

			if native {
				ancestors = ancestors[:len(ancestors)-1]
			}

			return more
		}

		walk(err, 0)
	}
}
