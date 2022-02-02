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

## LICENSE

This library is provided under BSD 3-Clause License, for more details see the LICENSE file.

