# oops

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

## Usage

### Fetch

Fetch the core library using:

```shell
go get go.sdls.io/oops/pkg/oops
```

### Define your errors

It's important to note, that with this library, errors can be defined globally. This allows:
- quickly identifying the errors of a package (usually grouped in a `errors.go` file, or at the top of each source file where the error is used mainly)
- using native Go error checks like `errors.Is` or `errors.As`

```go
var (
	ErrAuthMissing        = oops.Define("type", "auth", "code", "auth_missing", "status", 401)
	ErrAuthBadCredentials = oops.Define("type", "auth", "code", "auth_bad", "status", 401)
	ErrAuthExpired        = oops.Define("type", "auth", "code", "auth_expired", "status", 401)
)
```

In this example, we've defined three errors. In your projects/organization you may decide to build a helper function
that enforces the presence of certain fields, or even a more advanced Registry struct that holds all the errors.

### Yeet *your* errors

In `oops` we [`Yeet`](https://youtu.be/D8KxdXEBkhw) our errors. By default, when using `oops.Define()` you get
a `oops.ErrorDefined`, not something you should use directly as an `error`.

The reasons behind this is to have a verifiable error (`oops.ErrorDefined`) to which all other errors point. That's why
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
	ErrAuthMissing        = oops.Define("type", "auth", "code", "auth_missing", "status", 401)
	ErrAuthBadCredentials = oops.Define("type", "auth", "code", "auth_bad", "status", 401)
	ErrAuthExpired        = oops.Define("type", "auth", "code", "auth_expired", "status", 401)
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
		return oops.Explainf(err, "parsing token")
	}

	if tokenInfo.Expiry.Before(time.Now()) {
		return ErrAuthExpired.Yeetf("token expired at %s", tokenInfo.Expiry.Format(time.RFC3339))
	}
	
	return nil
}
```

As you can see above, we always call Yeet, with some explanation as to "the error occurred while doing X" and we can
even include formatted args.


### Wrap *their* errors

In the "Yeet *your* errors" example, we see how to return errors when _our_ code is the source of the error. But
sometimes you'll want to pass an error from the stdlib or a third party library. In that case, you can wrap the error
with `Wrap`.

### Custom Formatter

By default, the defined errors have a rudimentary string formatter that provides little (`Error.Explanation`) to no information regarding the error. Our recommended pattern is to have a dedicated package (be it locally in the project or as a organization level library) that wraps our top level functions calls such as `oops.Define` with typed arguments that represent **your** error handling params.

Eg: you might define `func Define(status int, code string) oops.ErrorDefined` and use that in your codebase with a formatter that then returns those `status` and `code` params the expected way.

## LICENSE

This library is provided under BSD 3-Clause License, for more details see the LICENSE file.

