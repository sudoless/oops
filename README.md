# oops

A custom developed library by SUDOLESS, tailored from our experience and needs.
Whilst using our first iteration of the bespoke error library (formally called `qer`), we noticed
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

### Fetch

Fetch the core library using:

```shell
go get go.sdls.io/oops@latest
```

### Define your errors

Errors are defined globally as `*ErrorDefinition` sentinels. Each definition has a **code** (identity string)
and optional builders for semantic tags, messages, tracing, inheritance, and custom formatting.

```go
var (
	ErrAuth        = oops.Define("auth").Causes(oops.CauseAuth).Actions(oops.ActionAuth)
	ErrAuthExpired = oops.Define("auth_expired").Causes(oops.CauseExpired).Actions(oops.ActionAuth).Inherits(ErrAuth)
	ErrNotFound    = oops.Define("not_found").Causes(oops.CauseNotFound).Message("resource not found")
	ErrValidation  = oops.Define("validation").Causes(oops.CauseValidation).Actions(oops.ActionFix)
)
```

Definition builders (all chainable):
- `.Causes(...)` — semantic tags describing **why** (e.g. `CauseAuth`, `CauseTimeout`, `CauseIO`)
- `.Actions(...)` — semantic tags describing **what to do** (e.g. `ActionRetry`, `ActionAbort`, `ActionFix`)
- `.Message(msg)` — public-facing message for clients
- `.Traced()` — enables stack trace capture on error creation
- `.Inherits(defs...)` — forms an inheritance chain so `errors.Is(child, parent)` works
- `.SetFormatter(fn)` — custom `func(*Error) string` for rendering

### Yeet *your* errors

In `oops` we [`Yeet`](https://youtu.be/D8KxdXEBkhw) our errors. `ErrorDefinition` is a sentinel — not
something you use directly as an `error`. Call `Yeet` or `Yeetf` to create a live `*Error` instance.

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

### Wrap *their* errors

When the error originates from the stdlib or a third party library, wrap it with `Wrap` or `Wrapf`
to attach it to one of your definitions.

```go
var ErrDatabase = oops.Define("database").Causes(oops.CauseIO).Actions(oops.ActionRetry)

func getUser(id string) (*User, error) {
	row := db.QueryRow("SELECT ...", id)
	if err := row.Scan(&user); err != nil {
		return nil, ErrDatabase.Wrapf(err, "scanning user %s", id)
	}
	return &user, nil
}
```

It is recommended you pair `oops` with a linter like [wrapcheck](https://github.com/tomarrell/wrapcheck).

### Semantic causes and actions

Every definition (and every live error) carries **cause** and **action** tags. Causes describe *why*
the error happened; actions describe *what the caller should do*.

Built-in causes: `CauseInternal`, `CauseNotFound`, `CauseAuth`, `CauseForbidden`, `CauseConflict`,
`CauseRateLimit`, `CauseTimeout`, `CauseBadRequest`, `CauseUnavailable`, `CauseBadGateway`,
`CauseExpired`, `CauseIO`, `CauseValidation`.

Built-in actions: `ActionRetry`, `ActionAbort`, `ActionFatal`, `ActionFix`, `ActionWait`,
`ActionAuth`, `ActionSkip`.

You can use your own string constants as well — `Cause` and `Action` are type aliases for `string`.

### Matching errors

The `Matcher` system lets you inspect errors by semantic tags without importing specific definitions.
This decouples error producers from consumers.

```go
// Retry any error that says it's retryable
if oops.Match(err, oops.ByCause(oops.CauseTimeout)) {
	return retry(fn)
}

// Combine matchers
isRetryable := oops.Any(
	oops.ByCause(oops.CauseTimeout),
	oops.ByAction(oops.ActionRetry),
)
if oops.Match(err, isRetryable) {
	return retry(fn)
}

// Negate
if oops.Match(err, oops.Not(oops.ByAction(oops.ActionFatal))) {
	// safe to retry
}
```

Available matchers: `ByCause`, `ByAction`, `ByCode`, `ByDefinition`.
Combinators: `All`, `Any`, `Not`.

### Definition inheritance

Definitions can inherit from other definitions. This allows broad matching via `errors.Is`:

```go
var (
	ErrAuth        = oops.Define("auth")
	ErrAuthExpired = oops.Define("auth_expired").Inherits(ErrAuth)
	ErrAuthRevoked = oops.Define("auth_revoked").Inherits(ErrAuth)
)

// All of these are true:
errors.Is(ErrAuthExpired.Yeet(), ErrAuth)     // true — via Inherits
errors.Is(ErrAuthExpired.Yeet(), ErrAuthExpired) // true — direct match
```

### Collecting batch errors

For operations that can produce multiple errors (e.g. validating a struct), use `Collect`:

```go
finish, add := ErrValidation.Collect()

if name == "" {
	add(ErrValidation.Yeetf("name is required"), "name")
}
if age < 0 {
	add(ErrValidation.Yeetf("age must be positive"), "age")
}

if err := finish(); err != nil {
	return err // wraps all collected errors under ErrValidation
}
```

### Package-level utilities

- `oops.Catch(err)` — extract `*Error` or wrap with `ErrUncaught`
- `oops.Assert(err)` — like `Catch` but returns `(*Error, bool)` where `bool` is `true` only when `err` was already an `*Error`
- `oops.Explainf(err, format, args...)` — append explanation to any error
- `oops.AddCause(err, causes...)` — append cause tags to any error
- `oops.Pathf(err, format, args...)` — set a path segment on any error
- `oops.As(err, target) (*Error, bool)` — find `*Error` in the unwrap chain matching a definition
- `oops.Nest(def, errs...)` — create a parent error wrapping multiple children

### Custom Formatter

The default formatter renders `code: message; explanation` when both are present, falling back
to `code: explanation`, `code: message`, or just `code` as applicable. For richer output, set a custom
`Formatter` on your definitions:

```go
var ErrAPI = oops.Define("api_error").
	Message("something went wrong").
	SetFormatter(func(err *oops.Error) string {
		return fmt.Sprintf("[%s] %s (%s)", err.Code(), err.Message(), err.Explanation())
	})
```

### Go compatible

`*Error` implements `Unwrap() []error` (Go 1.20+ multi-error), `Is`, and `Error`.
`*ErrorDefinition` implements `Is` and `Error`. You can use `errors.Is`, `errors.As`,
and `errors.Unwrap` with oops errors seamlessly.

## LICENSE

This library is provided under BSD 3-Clause License, for more details see the LICENSE file.
