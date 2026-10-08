// Package oops defines, creates, annotates, classifies and walks errors.
//
// # Lifecycle
//
// Define a sentinel once, at package level. A definition has a code, and
// builders add cause and action tags, a public message, tracing, inheritance
// and a custom formatter:
//
//	var ErrNotFound = oops.Define("not_found").Causes(oops.CauseNotFound).Message("resource not found")
//
// Create a live *Error from it: Yeet and Yeetf for your own errors, Wrap and
// Wrapf to attach somebody else's error. Annotate the *Error on the way up with
// Explainf, Set, AddCauses, SetActions, Pathf and Nest. Read it at the boundary
// with errors.Is, errors.As, As, Foreign, Native, the accessors, Walk and %+v.
//
// Every builder returns a new *ErrorDefinition and leaves its receiver
// unchanged. The sentinel is the pointer returned by the last call of the
// chain. A definition is never an error value: its Error method panics. Use
// Define(code).Inherits(parent) to derive a definition that matches its parent
// under errors.Is and As.
//
// # Foreign and Native
//
// Foreign returns an error as an *Error: an *Error is returned as is and any
// other error is wrapped with ErrForeign. Native reports whether an error is an
// *Error. Both inspect only the top-level value and never traverse wrapped
// errors. The package helpers Explainf, AddCauses, Pathf and Nest take any
// error and return error, with nil staying nil.
//
// A typed-nil *Error is treated as nil at every intake, and every *Error method
// accepts a nil receiver.
//
// # Ownership
//
// A live *Error has one owner and is mutated in place by its annotating
// methods. Give code that annotates the same error, or another goroutine, its
// own Clone. Clone is shallow: the wrapped errors, the trace and the field
// values are shared.
//
// # Reading the tree
//
// An *Error implements Unwrap() []error, so the Unwrap function of package
// errors returns nil on it. errors.Is, errors.As and As traverse wrapped
// errors, children and inherited definitions. (*Error).Is compares definitions
// only. Walk iterates the whole tree in pre-order, foreign errors included, and
// stops after 1024 nodes.
//
// Error returns the default rendering "code: message; explanation", dropping
// the empty parts, or the output of the definition's Formatter. Text of wrapped
// errors never appears in it. %v and %s print Error; %+v prints the tree: the
// path, causes, actions, fields and trace of each error, followed by the errors
// below it. A Formatter controls only the first line of each error in that
// tree.
//
// # Anti-patterns
//
//   - Wrapping an *Error in fmt.Errorf("%w") or errors.Join. Foreign and
//     Native read only the top-level value, so they classify the result as
//     foreign. Wrap with a definition instead: Def.Wrapf(err, ...).
//   - Returning or logging a definition without Yeet. Its Error method panics,
//     and fmt and slog print %!v(PANIC=Error method: ...).
//   - Declaring functions that return *Error instead of error. A nil *Error
//     stored in an error is a non-nil interface.
//   - Sharing an *Error across goroutines or call sites without Clone.
package oops
