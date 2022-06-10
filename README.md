# oops

[![CI](https://github.com/sudoless/oops/actions/workflows/ci.yaml/badge.svg)](https://github.com/sudoless/oops/actions/workflows/ci.yaml)

A custom developed library by SUDOLESS, tailored from our experience and needs.
Whilst using our first iteration of the bespoke error library (formally called `qer`), we noticed
a set of key issues and features that are important when dealing with errors in our services.

Before using any Go error library in a large or "future-proof codebase", please consider the proposed draft for
"[Go2 Errors](https://go.googlesource.com/proposal/+/master/design/go2draft.md)" and the
[feedback wiki](https://github.com/golang/go/wiki/Go2ErrorHandlingFeedback) for said draft. My personal opinion/feedback
is also [listed](https://gist.github.com/cpl/54ed073e20f03fb6f95257037d311420).

> This library is subject to breaking changes until it reaches v1.

## Why this library

### Perspective

Errors exist in two (equally real, equally important) perspectives. There is the "_CLIENT_" view and,
the "_SERVER/DEVELOPER/SYSOPS_" view.

Your clients/consumers should get relevant information as to **why**
the issue occurred, **what** are the results/consequences of the error and, **how** they can fix it.

You as the developer, system and/or admin, care as much about the above as the client, but you also care
about "unexpected" errors, and a more detailed view of the situation.

As such it is important to track  both internal information such as types, structs, lines of codes,
stack traces, etc for internal use and debugging, but it's also important to provide users with
"dumbed down" versions of the error.

### Human VS Machine

Another thing to consider when dealing with errors, is that errors will be handled by both machines and
humans. As such we provide easy to parse machine `code` and easier to understand human `explain`
(explanation) and `help` messages.

## Usage

### Fetch

Fetch the core library using `go get go.sdls.io/oops/pkg/oops`, other helper packages can be fetched like `go get go.sdls.io/oops/pkg/oops_json`.

### Define your errors

It's important to note, that with this library, errors should be defined globally. This allows:
- quickly identifying the errors of a package (usually grouped in a `errors.go` file, or at the top of each source file where the error is used mainly)
- using native Go error checks like `errors.Is` or `errors.As`
- creating groups of errors using `oops`

```go
var (
	ErrAuthMissing        = oops.Define().Type("auth").Code("auth_missing").StatusCode(401)
	ErrAuthBadCredentials = oops.Define().Type("auth").Code("auth_bad").StatusCode(401)
	ErrAuthExpired        = oops.Define().Type("auth").Code("auth_expired").StatusCode(401)
)
```

In this example, we've defined three errors. The same three errors could be defined as a group:

```go
var (
    errAuthGroup = oops.Define().Type("auth").StatusCode(401).Group().PrefixCode("auth_")

    ErrAuthMissing        = errAuthGroup.Code("missing")
    ErrAuthBadCredentials = errAuthGroup.Code("bad")
    ErrAuthExpired        = errAuthGroup.Code("expired")
)
```

In a group, errors will share the type, status code and other base fields. After calling `.Code` on a group, you can
modify the error in any way you want, without affecting the other errors in the group.

### Yeet *your* errors

In `oops` we [`Yeet`](https://youtu.be/D8KxdXEBkhw) our errors. By default, when using `oops.Define()` you get
a `*oops.errorDefined`, not something you should use directly (to make sure, calling `errorDefined.Error` will `panic`).

The reasons behind this is to have a verifiable error (`oops.errorDefined`) to which all other errors point. That's why
Yeet (or Wrap) must be called on a defined error when you want to _use it_.

```go
package example

import (
	"net/http"
	"strings"
	"time"

	"go.sdls.io/oops/pkg/oops"
)

var (
	ErrAuthMissing        = oops.Define().Type("auth").Code("auth_missing").StatusCode(401)
	ErrAuthBadCredentials = oops.Define().Type("auth").Code("auth_bad").StatusCode(401)
	ErrAuthExpired        = oops.Define().Type("auth").Code("auth_expired").StatusCode(401)
)

func validateAuth(r *http.Request) error {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ErrAuthMissing.Yeet("empty auth header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ErrAuthBadCredentials.Yeet("bad auth header, expected Bearer")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	
	// imagine we have an auth manager checking tokens
	expiredAt := authManager.CheckToken(token)

	if time.Now().After(expiredAt) {
		return ErrAuthExpired.Yeet("token expired at '%s'", expiredAt.Format(time.RFC3339))
	}

	return nil
}
```

As you can see above, we always call Yeet, with some explanation as it "why the error occurred" and even formatted args.


### Wrap *their* errors

In the "Yeet *your* errors" example, we see how to return errors when _our_ code is the source of the error. But
sometimes you'll want to pass an error from the stdlib or a third party library. In that case, you can wrap the error
with `Wrap`.


### Never do this

You must **NEVER** return a defined error directly.

```go
var ExampleErr = oops.Define().Type("example").Code("example_err").StatusCode(500)

func thinger1() error {
    // Do not do this ❌
    return ExampleErr
}

func thinger2() error {
    // Do not do this ❌
	return oops.Define().Type("example").Code("example_err").StatusCode(500)
}

func thinger3() error {
    // This is fine ✅
    return ExampleErr.Yeet("thinger 3")
}

func thinger4() error {
    // This is fine ✅
	err := errors.New("original error")
    return ExampleErr.Wrap(err, "thinger 4 with wrap")
}
```

### Handle errors


### Multiple reasons


### Stack traces


### Log errors



### HTTP integration


## LICENSE

This library is provided under BSD 3-Clause License, for more details see the LICENSE file.

