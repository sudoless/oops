# oops

A custom developed library by SUDOLESS, tailored from our experience and needs.
Whilst using our first iteration of the bespoke error library (formally called `qer`), we noticed
a set of key issues and features that are important when dealing with errors in our services.

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

## Library Concepts

### code

Each error has a human and machine (more targeted at machines, but some people fall in this category anyway)
readable three part code. Each code follows the convention `BLAME.NAMESPACE.REASON`. Do not take the terms
to mean their literal meaning. The three codes, and their purpose are mere _recommendations_, they are in
no way enforced by the library or thought police.

* `BLAME` covers a small set of "entities" that can be blamed for the error. The "entity" is relative
and not absolute. Meaning that in the right context `CLIENT` may be used in a `SERVER` "area".

* `NAMESPACE` is the context or "area" of code where the error occurred. It should narrow down certain errors.

* `REASON` is the most important of all three. `REASON` is a broad idea of "why" the error happened, also each `REASON`
is mapped to an HTTP Status Code, so it is important to pick the right `REASON` if you plan to serve the errors.

### define

The library uses two definitions of `error`s. One (`errorDefined` which is a private/internal struct) is meant to only
be used and created with the `oops.Define` function as follows:

```go
var (
    ErrDbConnString = oops.Define(oops.BlameClient, oops.NamespaceSetup, oops.ReasonRequestBad,
    "check your connection string and try again")
    ErrDbCreate = oops.Define(oops.BlameServer, oops.NamespaceSetup, oops.ReasonInternal)
    ErrDbConnect = oops.Define(oops.BlameServer, oops.NamespaceSetup, oops.ReasonConnection,
    "check your connection string and network configuration then try again")
)
```

When defined, the function takes an optional parameter that will end up as the "help" message for the error.

Errors should be "defined" at the top of the file (or in their own file). The structs returned by `oops.Define` MUST NOT
be used as Golang builtin `error` (as they will `panic` on any attempt to call `.Error()`). These errors are used as
parents/sources/generators for the other type (`oops.Error`).

### yeet

To generate a `oops.Error` you take one of your "defined" errors and call one of the methods (`Yeet`, `Wrap`) or any
of the other derivations. (e.g. `ErrDbConnect.Yeet()`)

`Yeet` will return a plain error, containing nothing more than the same Blame, Namespace and Reason as the parent
"defined" error.

### wrap

When you want to "wrap" another third party error, into an `oops.Error` then you can use for example
`return ErrDbConnect.Wrap(parentErr)`. This will keep track of the originating error, which can then be accessed with
the standard Go `errors.Unwrap` pattern, or checked via the `errors.Is` patterns.

### explain

An optional (but highly recommended) concept is the `oops.Explain` (or the helper functions `YeetExplain`, `WrapExplain`)
which will add your own flavour and meaning of "an explanation". What is an explanation is up to the reader.

The purpose of `Explain` is to call it on all (or nearly all) returns of an errors, this way you can build a
human-readable trace, which in the end may provide more help.

### multi

Each `oops.Error` has an optional method `Error.Multi()`, which can be used to return a string array as a more detailed
explanation of the error. We find this mandatory when returning validation errors, as to help people fix all errors
from the first try, instead of returning the validation errors one at a time (trial and error).

### stack trace

An optional method on the `oops.Error` si the `Error.Stack()` call which will populate the errors stack field ([]string)
with a line by line stack trace, of the functions which were called up to that point. This information is never made
available to the end user and should only be used by the developer.

## Compatability

The `oops.Error` can be used in any way any Go `error` can be used. It even provides support for checking if an error
occurred from a specific "defined" error using the `Error.Is` or `errors.Is` functions.

## Examples

```go
package main

import (
	"database/sql"
	"log"
	"os"
	"strings"

	"go.sdls.io/oops/pkg/oops"
)

var (
	ErrDbConnString = oops.Define(oops.BlameClient, oops.NamespaceSetup, oops.ReasonRequestBad,
		"check your connection string and try again")
	ErrDbCreate = oops.Define(oops.BlameServer, oops.NamespaceSetup, oops.ReasonInternal)
	ErrDbConnect = oops.Define(oops.BlameServer, oops.NamespaceSetup, oops.ReasonConnection,
		"check your connection string and network configuration then try again")
)

func main() {
	_, err := setup(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}

func setup(conn string) (*sql.DB, error) {
	if conn == "" {
		return nil, ErrDbConnString.YeetExplain("missing connection string")
	}
	if !strings.HasPrefix(conn, "fiz://") {
		return nil, ErrDbConnString.YeetExplain("connection string does not use fiz:// protocol")
    }

    db, err := sql.Open("foobar", conn)
	if err != nil {
		return nil, ErrDbCreate.Wrap(err)
	}

	if err = db.Ping(); err != nil {
		return nil, ErrDbConnect.WrapExplain(err, "ping db")
	}

	return db, nil
}
```






## LICENSE

This library is provided under BSD 3-Clause License, for more details see the LICENSE file.

