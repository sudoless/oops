# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```shell
# Run all tests
go test ./...

# Run a single test
go test ./pkg/oops/... -run TestName

# Run tests with race detector
go test -race ./...

# Lint (uses go-lint skill or golangci-lint directly)
golangci-lint run ./...
```

## Architecture

`oops` is a Go error library (`go.sdls.io/oops`) in the single package `pkg/oops`. There are no external dependencies.

### Two core types

- **`ErrorDefinition`** (`*ErrorDefinition` struct): A sentinel error definition created once at package level via `oops.Define(code)`. It holds a code string, semantic cause/action tags, a public-facing message, optional inheritance chain, tracing flag, and an optional `Formatter`. It implements `error` (returns `"code: message"`), `Is`, and creation methods (`Yeet`, `Yeetf`, `Wrap`, `Wrapf`, `Collect`).
- **`Error`** (`*Error` struct): A live error instance spawned from an `ErrorDefinition`. Implements `error`, `Unwrap() []error` (Go 1.20+ multi-error), `Is`. Holds explanation text, cause/action tags (copied from definition, mutable), wrapped errors, optional stack trace, path segment, and arbitrary fields map.

### Lifecycle

1. **Define** globally: `var ErrFoo = oops.Define("foo")` — optionally chain `.Causes()`, `.Actions()`, `.Message()`, `.Traced()`, `.Inherits()`, `.SetFormatter()`.
2. **Yeet** (own errors): `ErrFoo.Yeet()` / `ErrFoo.Yeetf(format, args...)` — creates an `*Error` with no wrapped parent.
3. **Wrap** (external errors): `ErrFoo.Wrap(err)` / `ErrFoo.Wrapf(err, format, args...)` — creates an `*Error` wrapping an external error.
4. **Collect** (batch errors): `finish, add := ErrFoo.Collect()` — accumulates errors via `add`, `finish()` returns nil if none added.

### Semantic classification

- **`Cause`** (`string` alias): Tags describing *why* an error occurred. Built-in constants: `CauseInternal`, `CauseNotFound`, `CauseAuth`, `CauseForbidden`, `CauseConflict`, `CauseRateLimit`, `CauseTimeout`, `CauseBadRequest`, `CauseUnavailable`, `CauseBadGateway`, `CauseExpired`, `CauseIO`, `CauseValidation`.
- **`Action`** (`string` alias): Tags describing *what the caller should do*. Built-in constants: `ActionRetry`, `ActionAbort`, `ActionFatal`, `ActionFix`, `ActionWait`, `ActionAuth`, `ActionSkip`.
- Causes/actions are set on definitions (inherited by errors) and can be mutated on live errors via `.AddCause()` / `.SetActions()`.

### Matcher system

`Matcher = func(*Error) bool` — composable predicates for inspecting errors by semantic tags without importing definitions.
- `Match(err, matcher)` — test any `error` against a `Matcher`.
- Pre-built: `ByCause(cause)`, `ByAction(action)`, `ByCode(code)`, `ByDefinition(def)`.
- Combinators: `All(...)`, `Any(...)`, `Not(m)`.

### Key package-level functions

- `oops.Catch(err)` — extract `*Error` or wrap with `ErrUncaught`.
- `oops.Assert(err)` — like `Catch` but returns `(*Error, bool)` flag.
- `oops.Explainf(err, format, args...)` — append explanation; wraps non-oops errors with `ErrUncaught`.
- `oops.AddCause(err, causes...)` — append cause tags; wraps non-oops errors with `ErrUncaught`.
- `oops.Pathf(err, format, args...)` — set path segment; wraps non-oops errors with `ErrUncaught`.
- `oops.As(err, target)` — traverse the unwrap chain (including multi-error) to find an `*Error` whose definition matches `target` (respects inheritance).
- `oops.Nest(def, errs...)` — create a parent error wrapping multiple children.
- `oops.Match(err, matcher)` — test error against a `Matcher` predicate.

### Error mutators (on `*Error`)

All return `*Error` for chaining:
- `.Explainf(format, args...)` — append explanation text.
- `.Set(key, value)` — store arbitrary field.
- `.AddCause(causes...)` — append cause tags.
- `.SetActions(actions...)` — replace action tags.
- `.WithPathf(format, args...)` — set path segment.
- `.Nest(err)` — add a single wrapped error.
- `.Append(errs...)` — add typed `*Error` children.

### Error accessors (on `*Error`)

- `.Definition()`, `.Code()`, `.Message()`, `.Explanation()` — identity and text.
- `.Causes()`, `.Actions()`, `.HasCause(c)`, `.HasAction(a)` — semantic tags.
- `.Fields()`, `.Get(key)` — arbitrary metadata.
- `.Path()`, `.PathArgs()` — path segment.
- `.Trace()` — stack frames (`[]string`).
- `.Unwrap()` — wrapped errors (`[]error`).

### Definition inheritance

`ErrorDefinition.Inherits(defs...)` forms a parent chain. `errors.Is` traverses the chain, enabling broad matching (e.g. `errors.Is(ErrAuthExpired.Yeet(), ErrAuth)` is true when `ErrAuthExpired.Inherits(ErrAuth)`).

### Preset errors

`preset.go` defines two built-in `ErrorDefinition` values:
- `ErrUncaught` — auto-wraps non-oops errors (traced, `CauseInternal`, `ActionAbort`).
- `ErrTODO` — placeholder for unimplemented paths (traced, `CauseInternal`, `ActionAbort`, message "not implemented").

### Stack tracing

`internal/unsafe/stack.go` captures runtime stack frames via `runtime.Caller`. Only activated when `ErrorDefinition.traced == true` (set via `.Traced()` builder or preset errors).

### Go compatibility

`*Error` implements `Unwrap() []error` (Go 1.20+ multi-error unwrap), `Is`, and `Error`. `*ErrorDefinition` implements `Error` and `Is`. Both work with `errors.Is` / `errors.As` / `errors.Unwrap`. The `oops.As` function provides definition-aware traversal across multi-error trees.
