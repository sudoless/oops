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

- **`ErrorDefined`** (`*errorDefined`): A sentinel error definition created once at package level via `oops.Define(key, value, ...)`. Never used directly as `error` — calling `.Error()` panics. It holds key-value props and an optional `Formatter`.
- **`Error`** (`*errorImpl`): A live error instance spawned from an `ErrorDefined`. Implements `error`, `Unwrap`, `Is`, `As`. Holds explanation text, key-value props (copied from source), optional stack trace, optional nested errors, and an optional path string.

### Lifecycle

1. **Define** globally: `var ErrFoo = oops.Define("key", "val")` — optionally chain `.Trace()`, `.Set()`, `.Formatter()` builders.
2. **Yeet** (own errors): `ErrFoo.Yeet()` / `ErrFoo.Yeetf(format, args...)` — creates an `errorImpl` with no parent.
3. **Wrap** (external errors): `ErrFoo.Wrap(err)` / `ErrFoo.Wrapf(err, format, args...)` — creates an `errorImpl` with a parent.
4. **Collect** (batch errors): `finish, addf := ErrFoo.Collect()` — accumulates errors into `Error.Nested`, returns nil if none added.

### Key package-level functions

- `oops.Explainf(err, format, args...)` — appends explanation to an existing `Error`; wraps non-`Error` values with `ErrUncaught`.
- `oops.As(err, target)` — traverses the unwrap chain to find an `Error` whose `.Source()` matches `target`.
- `oops.Nest(source, nested...)` — creates a parent error with nested children.
- `oops.NestedAs` / `oops.NestedIs` — search through `Error.Nested` recursively (does not traverse unwrap chain).

### Preset errors

`preset.go` defines two built-in `ErrorDefined` values:
- `ErrUncaught` — auto-wraps non-`oops` errors (traced).
- `ErrTODO` — placeholder for unimplemented paths (traced).

### Stack tracing

`internal/unsafe/stack.go` captures runtime stack frames via `runtime.Caller`. Only activated when `errorDefined.traced == true` (set via `.Trace()` builder or preset errors).

### Go compatibility

`errorImpl` implements `Is` and `As` to work with `errors.Is` / `errors.As`. `errors.Is` on an `Error` checks `Source()` equality against an `ErrorDefined` or another `Error`. `errors.As` with `*Error` as target sets the pointer to the concrete `errorImpl`.