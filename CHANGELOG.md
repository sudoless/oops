# Changelog

## v2.0.0 (Unreleased)

### Symbols

| v1.0.1 | v2 | Notes |
|---|---|---|
| import `go.sdls.io/oops/pkg/oops` | import `go.sdls.io/oops/v2/pkg/oops` | module path `go.sdls.io/oops/v2` |
| interface `ErrorDefined` | `*ErrorDefinition` | a struct pointer; it is the sentinel's identity |
| interface `Error` | `*Error` | a struct pointer; use it as `*oops.Error` |
| `Define(props ...any)` | `Define(code string)` | tags and message come from builders |
| `.Trace()` | `.Traced()` | |
| `.Set(key, value)` on a definition | removed | set fields on the error: `(*Error).Set` |
| `.Formatter(f)` | `.Formatter(f)` | `Formatter` is now `func(*Error) string` |
| n/a | `.Causes(...)`, `.Actions(...)`, `.Message(msg)`, `.Inherits(defs...)` | new definition builders; `Cause*` and `Action*` constants |
| definition props | `Set`, `Get` on errors | definitions no longer carry key-value props |
| `GetAll()` | `Fields()` | |
| `Source()` | `Definition()` | returns `*ErrorDefinition` |
| `Append(errs ...Error)` | `Nest(errs ...error)` | nil and typed-nil `*Error` values are skipped |
| `Nested()` | removed | children are returned by `Unwrap()` and visited by `Walk` |
| `PathSetf(path, args...)` | `Pathf(format, args...)` | |
| `Explainf(format, args...)` (no result) | `Explainf(format, args...) *Error` | |
| `Unwrap() error` | `Unwrap() []error` | `errors.Unwrap` returns nil; use `errors.Is`, `errors.As`, `oops.As` or `oops.Walk` |
| `As(other any) bool` | removed | `errors.As` finds the first `*Error` in the tree |
| `var e oops.Error; errors.As(err, &e)` | `var e *oops.Error; errors.As(err, &e)` | the old form compiles but panics at run time; `go vet` reports it |
| `def.Is(err)` | `errors.Is(err, def)` | `(*ErrorDefinition).Is` matches definitions only |
| n/a | `(*ErrorDefinition).Code()` | returns the definition's code |
| n/a | `Code`, `Message`, `Causes`, `Actions`, `AddCauses`, `SetActions`, `HasCause`, `HasAction`, `Clone`, `Format` | new `*Error` methods |
| `MustAny(err) Error` | `Foreign(err) *Error` | wraps with `ErrForeign`; typed-nil `*Error` gives nil |
| `AssertAny(err) (Error, bool)` | `Native(err) (*Error, bool)` | `(nil, false)` for nil, typed-nil `*Error` and non-oops errors |
| `Explainf(err, format, args...) Error` | `Explainf(err, format, args...) error` | wraps a non-oops error with `ErrForeign`; nil stays nil; chain with `oops.Foreign(err).Explainf(...)` |
| n/a | `AddCauses(err, causes...) error`, `Pathf(err, format, args...) error` | same contract as `Explainf` |
| `Nest(source ErrorDefined, nested ...Error) Error` | `Nest(def *ErrorDefinition, errs ...error) error` | nil when `def` is nil or every error is nil |
| `As(err, ErrorDefined) (Error, bool)` | `As(err, *ErrorDefinition) (*Error, bool)` | searches wrapped errors and children, and matches `Inherits` parents |
| `NestedAs`, `NestedIs` | removed | children take part in `errors.Is`, `errors.As`, `oops.As` and `Walk` |
| n/a | `Walk(err) iter.Seq2[int, error]` | pre-order iterator over the error tree |
| `ErrUncaught` | `ErrForeign` | code `foreign`; renders `foreign` or `foreign: explanation` |
| `ErrTODO` | `ErrTODO` | renders `todo: not implemented` or `todo: not implemented; explanation` |
| `NilErr` | removed | use `nil` |
| `ErrorCollectorFinish = func() Error` | `CollectorFinish = func() error` | |
| `ErrorCollectorAdd` | `CollectorAdd` | same signature |

### Behaviour

| Area | v1.0.1 | v2 |
|---|---|---|
| Builders | mutate the definition and return it | copy-on-write: each builder returns a new `*ErrorDefinition` with a new identity, so a builder call on its own line has no effect and errors of a derived definition do not match its receiver |
| `Error()` | the explanation, or `oops.Error` without one | `code: message; explanation`, `code: explanation`, `code: message` or `code`, or the definition's `Formatter`; wrapped errors' text never appears |
| `%+v` | same as `%v` | multi-line tree: path, causes, actions, fields, trace and every error below, including foreign text; a `Formatter` controls only the first line of each error |
| `%#v` | Go-syntax dump of the struct | `Error()` as a quoted string, like `%q` |
| `(*Error).Is` | definition match, then `errors.Is` on the parent | compares definitions only, including `Inherits` parents; `errors.Is` traverses the tree |
| `(*ErrorDefinition).Is` | matches an `Error` created by the definition | matches only definitions: itself and the ones it inherits; use `errors.Is(err, def)` for an error |
| `errors.Is`, `errors.As`, `oops.As` | follow the parent chain | visit wrapped errors and children; `oops.As` stops at an `*Error` that appears below itself |
| `oops.As`, `%+v` | no node cap | stop after 1024 nodes, foreign ones included; beyond that `oops.As` can disagree with `errors.Is` and `errors.As`, and `%+v` ends with a line saying it was truncated |
| `Explainf`, `Pathf` formats | a format without arguments is used as is | every format goes through `fmt`: write `%%` for a literal `%` |
| `Pathf` arguments | `PathSetf` stored the caller's args slice | the arguments are copied: the caller may reuse the slice |
| `Pathf` | sets the path; arguments from a previous call stay in `PathArgs` | the path is a label: it overwrites, an empty format clears it, and `PathArgs` is reset on each call |
| Typed-nil `*Error` | counted as an error by `Wrap`, `Nest` and `Collect` | treated as nil by `Wrap`, `Wrapf`, `Nest`, `Collect`, `Foreign`, `Native`, `As` and `Is`; every `*Error` method accepts a nil receiver |
| `Collect` add | sets the path of the added error, even when empty; panics on an error that is not an `Error` or `ErrorDefined` | an empty path leaves the added error's path alone; a non-oops error is wrapped with `ErrForeign` |
| `Collect` finish | the result shares the collector's slice | the result has its own copy |
| Traces | start inside the library for some entry points; uncapped | start at the code that called the library; capped at 32 program counters |
| Copying an `Error` value | not possible: the value type was unexported | `go vet` reports `c := *e`; use `Clone()` |
| A definition passed where an error is expected (`Collect` add, `Wrap`, `Wrapf`, `Nest`, `Foreign`) | `Collect` add called `Yeet()` on it | an anti-pattern: it is treated as a foreign error, and rendering it panics because the definition's `Error()` panics; create an error with `Yeet` first |

## v1.0.1 Released (2026-03-05)

* Add CLAUDE
* Bug fix — `oops.As` dead code path (`Unwarp` typo)
* Add tests for coverage

## v1.0.0 Released (2026-01-29)

* Includes `oops` as is, releasing the first `v1` tag.
* Migration from `v0` depends on how old the `v0.x` is:
  * `v0.12+` should be enough to just go update (point to v1.0.0)
  * `v0.10+` might require a little refactoring, but overall possible
  * `pre-v0.10` requires major refactoring
  * `pre-v0.3` complete refactoring
