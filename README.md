# oops

A custom developed library by SUDOLESS, tailored from our experience and needs.
Whilst using our first iteration of the bespoke error library (formerly called `qer`), we noticed
a set of key issues and features that are important when dealing with errors in our services.

Before using any Go error library in a large or "future-proof codebase", please consider the proposed draft for
"[Go2 Errors](https://go.googlesource.com/proposal/+/master/design/go2draft.md)" and the
[feedback wiki](https://github.com/golang/go/wiki/Go2ErrorHandlingFeedback) for said draft. My personal opinion/feedback
is also [listed](https://gist.github.com/cpl/54ed073e20f03fb6f95257037d311420).

## Why this library

### Perspective

Errors exist in two (equally real, equally important) perspectives. There is the "_CLIENT_" view and,
the "_SERVER/DEVELOPER/SYSOPS_" view.

Your clients/consumers should get relevant information as to **why**
the issue occurred, **what** are the results/consequences of the error and, **how** they can fix it.

You as the developer, system and/or admin, care as much about the above as the client, but you also care
about "unexpected" errors, and a more detailed view of the situation.

As such it is important to track both internal information such as types, structs, lines of codes,
stack traces, etc for internal use and debugging, but it's also important to provide users with
"dumbed down" versions of the error.

## Usage

### Install

```shell
go get go.sdls.io/oops/v2
```

```go
import "go.sdls.io/oops/v2/pkg/oops"
```

The module has no dependencies. Migrating from v1: see the [CHANGELOG](CHANGELOG.md).

### Lifecycle

1. **Define** a sentinel once, at package level: `oops.Define(code)` plus builders.
2. **Create** a live error from it: `Yeet`/`Yeetf` for your own errors, `Wrap`/`Wrapf` for somebody else's.
3. **Annotate** on the way up: `Explainf`, `Set`, `AddCauses`, `SetActions`, `Pathf`, `Nest`.
4. **Read** at the boundary: `errors.Is`, `oops.As`, `oops.Foreign`, `oops.Native`, the accessors, `oops.Walk`, `%+v`.

There are two types. `*oops.ErrorDefinition` is the sentinel: its identity is its pointer, and it is never an
error value itself. `*oops.Error` is the live error created from it.

### Define your errors

Definitions are package-level `*ErrorDefinition` sentinels. Each has a **code** (identity string) and optional
builders for semantic tags, a public message, tracing, inheritance and custom formatting.

```go
var (
	ErrAuth        = oops.Define("auth").Causes(oops.CauseAuth).Actions(oops.ActionAuth)
	ErrAuthExpired = oops.Define("auth_expired").Causes(oops.CauseExpired).Actions(oops.ActionAuth).Inherits(ErrAuth)
	ErrNotFound    = oops.Define("not_found").Causes(oops.CauseNotFound).Message("resource not found")
	ErrValidation  = oops.Define("validation").Causes(oops.CauseValidation).Actions(oops.ActionFix)
)
```

Builders (all chainable):
- `.Causes(...)`: semantic tags describing **why** (`CauseAuth`, `CauseTimeout`, `CauseIO`, ...)
- `.Actions(...)`: semantic tags describing **what to do** (`ActionRetry`, `ActionAbort`, `ActionFix`, ...)
- `.Message(msg)`: public-facing message for clients
- `.Traced()`: capture a stack trace when an error is created
- `.Inherits(defs...)`: `errors.Is(err, parent)` also matches errors of this definition
- `.Formatter(fn)`: custom `func(*Error) string` for `Error()`

Every builder returns a **new** definition and leaves its receiver unchanged. The sentinel is the pointer returned
by the last call in the chain, so always assign the whole chain:

```go
var (
	ErrOK  = oops.Define("ok").Message("fine")
	_      = ErrOK.Message("changed") // result discarded: ErrOK is unchanged
	ErrNew = ErrOK.Message("changed") // a different sentinel; errors of ErrNew do not match ErrOK
)
```

### Yeet *your* errors

In `oops` we [`Yeet`](https://youtu.be/D8KxdXEBkhw) our errors. Call `Yeet` or `Yeetf` on a definition to create a
live `*Error`.

```go
var (
	ErrAuthMissing        = oops.Define("auth_missing").Causes(oops.CauseAuth).Actions(oops.ActionAuth)
	ErrAuthBadCredentials = oops.Define("auth_bad").Causes(oops.CauseAuth).Actions(oops.ActionAuth)
	ErrAuthExpired        = oops.Define("auth_expired").Causes(oops.CauseExpired).Actions(oops.ActionAuth)
)

func validateAuth(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ErrAuthMissing.Yeetf("empty auth header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ErrAuthBadCredentials.Yeetf("bad auth header [%s], expected Bearer", authHeader)
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	tokenInfo, err := authManager.ParseToken(token)
	if err != nil {
		return ErrAuthBadCredentials.Wrapf(err, "parsing token")
	}

	if tokenInfo.Expiry.Before(time.Now()) {
		return ErrAuthExpired.Yeetf("token expired at %s", tokenInfo.Expiry.Format(time.RFC3339))
	}

	return nil
}
```

`Error()` renders `code: message; explanation`, dropping the parts that are empty: `auth_missing: empty auth header`,
`not_found: resource not found; user 42`, `not_found`. A definition's `Formatter` replaces this format.

### Wrap *their* errors

When the error comes from the standard library or a third party, attach it to one of your definitions with `Wrap`
or `Wrapf`. The original stays reachable through `errors.Is`, `errors.As`, `oops.As` and `oops.Walk`.

```go
var ErrDatabase = oops.Define("database").Causes(oops.CauseIO).Actions(oops.ActionRetry)

func getUser(id string) (*User, error) {
	var user User
	row := db.QueryRow("SELECT ...", id)
	if err := row.Scan(&user); err != nil {
		return nil, ErrDatabase.Wrapf(err, "scanning user %s", id)
	}
	return &user, nil
}
```

The text of the wrapped error does not appear in `Error()`: it is available through `%+v` and `oops.Walk`.
It is recommended you pair `oops` with a linter like [wrapcheck](https://github.com/tomarrell/wrapcheck).

### Causes and actions

Every definition and every live error carries **cause** and **action** tags. Causes describe *why* the error
happened; actions describe *what the caller should do*. A new error starts with copies of its definition's tags.

Built-in causes: `CauseInternal`, `CauseNotFound`, `CauseAuth`, `CauseForbidden`, `CauseConflict`,
`CauseRateLimit`, `CauseTimeout`, `CauseBadRequest`, `CauseUnavailable`, `CauseBadGateway`,
`CauseExpired`, `CauseIO`, `CauseValidation`.

Built-in actions: `ActionRetry`, `ActionAbort`, `ActionFatal`, `ActionFix`, `ActionWait`,
`ActionAuth`, `ActionSkip`.

`Cause` and `Action` are aliases for `string`, so your own string constants work too.

```go
func shouldRetry(err error) bool {
	e, ok := oops.Native(err)
	return ok && (e.HasAction(oops.ActionRetry) || e.HasCause(oops.CauseTimeout))
}

func tag(err *oops.Error) {
	err.AddCauses(oops.CauseTimeout) // appends
	err.SetActions(oops.ActionWait)  // replaces
}
```

### Annotate

Mutators change the error in place and return it for chaining:

```go
func lookup(id int, tenant string) error {
	return ErrNotFound.Yeet().
		Explainf("user %d", id).
		Set("user_id", id).
		Set("tenant", tenant).
		Pathf("users[%d]", id)
}
```

- `Explainf(format, args...)` appends to the explanation (joined with `, `). The format always goes through `fmt`:
  write `%%` for a literal `%`.
- `Set(key, value)` / `Get(key)` / `Fields()` store and read key-value fields.
- `Pathf(format, args...)` sets a **label** for where the error belongs (a field name, an index). Setting it again
  replaces it; an empty format clears it.
- `Nest(errs...)` adds child errors. Nil and typed-nil `*Error` values are skipped.

### Package helpers

`Explainf`, `AddCauses`, `Pathf` and `Nest` also exist as package functions. They take any `error`, wrap a non-oops
error with `ErrForeign` first, and **return `error`**: nil stays nil.

```go
func load(id string) error {
	err := store.Get(id)
	return oops.Explainf(err, "loading %s", id)
}

func validate(name, email string) error {
	return oops.Nest(ErrValidation, checkName(name), checkEmail(email)) // nil if both are nil
}
```

To keep chaining on a foreign error, go through `oops.Foreign`:

```go
func loadPath(id string) error {
	err := store.Get(id)
	if err == nil {
		return nil
	}
	return oops.Foreign(err).Pathf("store[%s]", id).AddCauses(oops.CauseIO)
}
```

### Foreign and Native

`oops.Foreign(err)` returns an `*Error`: the error itself when it already is one, otherwise a new `ErrForeign`
error wrapping it (code `foreign`, traced, cause `internal`, action `abort`). It returns nil for a nil input.

`oops.Native(err)` is a pure type check: `(*Error, true)` for an `*Error`, `(nil, false)` for anything else.

Both inspect **only the top-level value**. They never look inside wrapped errors.

```go
func classify(err error) string {
	if e, ok := oops.Native(err); ok {
		return "native:" + e.Code()
	}
	return "foreign:" + oops.Foreign(err).Code()
}
```

Because `Foreign` returns `*Error`, check it for nil before returning it as an `error` (see Anti-patterns). The
package helpers do that check for you.

### Is and As

`errors.Is` and `errors.As` traverse the whole tree: wrapped errors, nested children, and `Inherits` parents.
`oops.As(err, def)` does the same walk and returns the first `*Error` whose definition is, or inherits, `def`.

```go
var (
	ErrAuth        = oops.Define("auth")
	ErrAuthExpired = oops.Define("auth_expired").Inherits(ErrAuth)
)

func match(err error) {
	errors.Is(err, ErrAuth)        // true for ErrAuth and ErrAuthExpired errors
	errors.Is(err, ErrAuthExpired) // true for ErrAuthExpired errors only

	if e, ok := oops.As(err, ErrAuth); ok {
		fmt.Println(e.Code(), e.Causes(), e.Explanation())
	}

	var oe *oops.Error
	if errors.As(err, &oe) { // first *oops.Error in the tree
		fmt.Println(oe.Code())
	}
}
```

`(*Error).Is` itself compares definitions only and never looks at wrapped errors; use `errors.Is` to traverse.
`*Error` implements `Unwrap() []error`, so the `Unwrap` function of package `errors` returns nil on it: use
`errors.Is`, `errors.As`, `oops.As` or `oops.Walk`.

### Collect

For operations that produce several errors (validating a struct), use `Collect`. Each added error may carry a
path label; non-oops errors are wrapped with `ErrForeign`; nils are skipped.

```go
func validateUser(name string, age int) error {
	finish, add := ErrValidation.Collect()

	if name == "" {
		add(ErrValidation.Yeetf("name is required"), "name")
	}
	if age < 0 {
		add(ErrValidation.Yeetf("age must be positive"), "age")
	}
	add(checkEmail(name), "") // an empty path leaves the error's own path alone

	return finish() // nil when nothing was added
}
```

`finish()` returns a new `ErrValidation` error whose children are the added errors, in order. The add and finish
functions are not safe for concurrent use.

### Walk

`oops.Walk(err)` iterates `err` and every error below it, in pre-order, with the depth of each node. It follows
`Unwrap() error` and `Unwrap() []error`, so foreign errors are included.

```go
func dump(err error) {
	for depth, node := range oops.Walk(err) {
		fmt.Printf("%*s%v\n", depth*2, "", node)
	}
}
```

A cycle through an `*Error` is cut where it repeats, and the walk stops after 1024 nodes.

### Printing with %+v

`%v` and `%s` print `Error()`; `%q` prints it quoted. `%+v` prints the whole tree: the error, its path, causes,
actions, fields and trace, then every error below it, indented, including the text of foreign errors.

```go
var (
	ErrService  = oops.Define("service").Message("could not load user")
	ErrDatabase = oops.Define("database").Causes(oops.CauseIO).Actions(oops.ActionRetry)
)

func show() {
	_, err := os.Open("/nonexistent")
	dbErr := ErrDatabase.Wrapf(err, "scanning user %d", 42).Set("user_id", 42).Pathf("users[%d]", 42)
	top := ErrService.Wrap(dbErr).Explainf("request %s", "r1")

	fmt.Printf("%v\n\n%+v\n", top, top)
}
```

```text
service: could not load user; request r1

service: could not load user; request r1
  └ database: scanning user 42
    path: users[42]
    causes: io
    actions: retry
    fields: user_id=42
    └ open /nonexistent: no such file or directory
```

`log/slog`'s text handler formats values with `%+v`, so it logs this full tree, including foreign text and trace
frames. The JSON handler uses `Error()`.

### Traces

A definition built with `.Traced()` captures a stack trace for each error it creates. The first frame is the code
that called `Yeet`, `Wrap` or the package helper. At most 32 program counters are captured. Read them with
`err.Trace()`; `%+v` prints them. `ErrForeign` and `ErrTODO` are traced.

### Custom Formatter

The default rendering is `code: message; explanation`, dropping the parts that are empty. A `Formatter` replaces it
for the errors of one definition:

```go
var ErrAPI = oops.Define("api_error").
	Message("something went wrong").
	Formatter(func(err *oops.Error) string {
		return fmt.Sprintf("[%s] %s (%s)", err.Code(), err.Message(), err.Explanation())
	})
```

The `Formatter` controls `Error()`, and so the first line of each error in `%+v`. The path, causes, actions,
fields, trace and child lines are always rendered by `oops`.

### Ownership and Clone

A live `*Error` has one owner and is mutated in place by `Explainf`, `Set`, `AddCauses`, `SetActions`, `Pathf`
and `Nest`, and by the package helpers when given an `*Error`. Do not copy the struct value. To hand the same
error to code that annotates it, or to another goroutine, give that code its own `Clone`:

```go
func fanOut(base *oops.Error, shards []string) []error {
	errs := make([]error, 0, len(shards))
	for _, shard := range shards {
		errs = append(errs, base.Clone().Explainf("shard %s", shard).Set("shard", shard))
	}
	return errs
}
```

`Clone` is a shallow copy: the copy owns its explanation, causes, actions, fields map, path, path args and list of
wrapped errors. The wrapped errors, the trace and the field values are shared.

### Go compatible

`*Error` implements `Error`, `Unwrap() []error` (Go 1.20+ multi-error), `Is` and `fmt.Formatter`.
`*ErrorDefinition` implements `Is` so that it can be an `errors.Is` target; its `Error()` panics, because a
definition is not an error value.

## Anti-patterns

**Wrapping an oops error in `fmt.Errorf("%w")` or `errors.Join`.** `errors.Is`, `errors.As`, `oops.As` and
`oops.Walk` still find the error, but `oops.Foreign` and `oops.Native` read only the top-level value: they see a
foreign error. `Foreign` then wraps it in a second, `foreign` error, and the original's tags are no longer at the
top.

```go
func wrongWrap(err *oops.Error) {
	wrapped := fmt.Errorf("loading: %w", err)

	if _, ok := oops.Native(wrapped); !ok { // top level is a *fmt.wrapError
		fmt.Println(oops.Foreign(wrapped).Code()) // "foreign", not err.Code()
	}
}

func rightWrap(err *oops.Error) error {
	return ErrService.Wrapf(err, "loading") // an oops error on top
}
```

**Returning or logging a definition without `Yeet`.** A definition satisfies `error` so it can be an `errors.Is`
target, but its `Error()` panics. `fmt` and `slog` recover the panic and print `%!v(PANIC=Error method: ...)`. Passing a definition where an error is
expected (`Collect` add, `Wrap`, `Wrapf`, `Nest`, `Foreign`) has the same effect: it is treated as a foreign error, and
rendering it panics.

```go
func wrongReturn() error {
	return ErrService // compiles; calling Error() panics
}

func rightReturn() error {
	return ErrService.Yeet()
}
```

**Returning `*oops.Error` instead of `error`.** A nil `*oops.Error` stored in an `error` is a non-nil interface, so
`err != nil` is true for it.

```go
func check(ok bool) *oops.Error {
	if ok {
		return nil
	}
	return ErrService.Yeet()
}

func wrongRun() error {
	return check(true) // non-nil error interface holding a nil *oops.Error
}

func rightRun() error {
	if e := check(true); e != nil {
		return e
	}
	return nil
}
```

The package helpers (`Explainf`, `AddCauses`, `Pathf`, `Nest`) return `error` for this reason; declare your own
functions to return `error` too.

**Sharing an `*oops.Error` across goroutines or call sites without `Clone`.** Annotating mutates the error in place,
so two owners overwrite each other's explanation, fields and path, and concurrent use is a data race.

```go
func wrongShare(shared *oops.Error) {
	go func() { shared.Explainf("worker a") }()
	go func() { shared.Explainf("worker b") }()
}

func rightShare(shared *oops.Error) {
	go func() { shared.Clone().Explainf("worker a") }()
	go func() { shared.Clone().Explainf("worker b") }()
}
```

## Rosetta Stone

How the same job reads in `oops`, the standard library, and three other error libraries. `n/a` means the library
has no equivalent.

| | oops | Go stdlib (`errors`, `fmt`) | `github.com/pkg/errors` | `github.com/hashicorp/go-multierror` | `github.com/cockroachdb/errors` |
|---|---|---|---|---|---|
| Define a sentinel | `var ErrX = oops.Define("x")` | `var ErrX = errors.New("x")` | `var ErrX = errors.New("x")` [^stack-init] | n/a (use stdlib) | `var ErrX = errors.New("x")` [^stack-init] |
| Create | `ErrX.Yeetf("load %d", id)` | `fmt.Errorf("load %d", id)` | `errors.Errorf("load %d", id)` (no `%w`) | n/a | `errors.Newf("load %d", id)` |
| Wrap with context | `ErrX.Wrapf(err, "load %d", id)` | `fmt.Errorf("load %d: %w", id, err)` | `errors.Wrapf(err, "load %d", id)` | n/a (`multierror.Prefix` only prefixes) | `errors.Wrapf(err, "load %d", id)` |
| Check identity | `errors.Is(err, ErrX)` (also matches `Inherits` parents) | `errors.Is(err, ErrX)` | `errors.Is(err, ErrX)` or `errors.Cause(err) == ErrX` | `errors.Is(err, ErrX)` | `errors.Is(err, ErrX)` [^crdb-is] |
| Extract a typed value | `oops.As(err, ErrX)` or `errors.As(err, &e)` | `errors.As(err, &t)`; Go 1.26+: `errors.AsType[T](err)` | `errors.As(err, &t)` | `errors.As(err, &t)` | `errors.As(err, &t)` (also follows `Cause()`) |
| Attach fields | `err.Set("id", 42)`, `err.Get("id")` | n/a | n/a | n/a | n/a (closest: `errors.WithContextTags(err, ctx)`) |
| Classify | `Causes(...)`, `Actions(...)`; `err.HasCause(c)` | n/a beyond sentinels and custom types | n/a | n/a | `WithDomain`, `WithTelemetry`, `WithHint`, `WithAssertionFailure`, `UnimplementedError` |
| Multiple errors | `ErrX.Collect()`, `oops.Nest(ErrX, e1, e2)` | `errors.Join(e1, e2)` | n/a | `multierror.Append(result, err)`, `multierror.Group` | `errors.Join(e1, e2)`, `errors.CombineErrors(e1, e2)` [^crdb-combine] |
| Stack trace | `.Traced()` on the definition; `err.Trace()` | none | automatic on `New`, `Errorf`, `Wrap`, `Wrapf`, `WithStack` | none | automatic on `New`, `Newf`, `Wrap`, `Wrapf`, `Join`, `WithStack` |
| Render the full chain | `fmt.Printf("%+v", err)`, `oops.Walk(err)` | `%+v` prints the same as `%v` | `%+v`: each message followed by its stack frames | `Error()` (and `%+v`) prints a bullet list | `%+v`: numbered `Wraps:` tree with stacks, hints and details |

[^stack-init]: The stack is captured when the package initialises, so a sentinel carries an init-time stack.
[^crdb-is]: `Is` also matches after a network encode/decode and through `errors.Mark`.
[^crdb-combine]: `CombineErrors` adds the second error as secondary: it is printed, but invisible to `errors.Is`.

Versions checked: Go 1.27.1, `github.com/pkg/errors` v0.9.1, `github.com/hashicorp/go-multierror` v1.1.1,
`github.com/cockroachdb/errors` v1.14.0. `errors.AsType` needs Go 1.26 or newer. `pkg/errors` is in maintenance
mode, with its last release in 2020. The `go-multierror` repository README on its main branch recommends
`errors.Join` for new projects; the v1.1.1 release does not contain that note.

## LICENSE

This library is provided under BSD 3-Clause License, for more details see the LICENSE file.
