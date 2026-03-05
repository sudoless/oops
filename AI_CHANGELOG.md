# AI Changelog

Changes made by Claude Code (AI assistant), dated for traceability.

---

## 2026-03-05

### Bug fix — `oops.As` dead code path (`Unwarp` typo)

**File:** `pkg/oops/functions.go`

`asErr` contained a type-switch case matching `interface{ Unwarp() error }` — a misspelling of `Unwrap`. Because no real type implements this misspelled method, the single-unwrap-chain branch was permanently dead. This meant `oops.As` silently failed to traverse any error wrapped with `fmt.Errorf("%w", oopsErr)`. The multi-error branch (`Unwrap() []error`, used by `errors.Join`) was correct and unaffected.

Fixed by correcting `Unwarp` → `Unwrap` on both the case and the call.

### New tests

Coverage increased from **68% → 86.8%**.

Tests were added directly into the relevant existing test files:

| File | Tests added |
|---|---|
| `functions_test.go` | `TestAs_fmtErrorfWrap`, `TestAs_fmtErrorfWrap_notFound` (regression for the bug above), `TestAssertAny`, `TestMustAny`, `TestNest`, `TestNestedIs` |
| `error_impl_test.go` | `TestError_Set`, `TestError_Get_missing`, `TestError_Append`, `TestError_Path`, `TestError_Explainf_emptyFormat` |
| `error_go_test.go` | `TestError_As_definedTargetPanics` |
| `define_impl_test.go` | `TestDefine_oddArgsPanic` |

### Remaining known gaps (intentionally not covered)

- `error_go.go:Error` — the `err.source.formatter == nil` branch is unreachable through the public API (all code paths that produce an `errorDefined` always set a formatter).
- `define.go:defaultFormatter` — the `err == nil` guard inside the formatter is unreachable because `errorImpl.Error()` returns early before calling the formatter on a nil receiver.
- `functions.go:asErr` — the `Unwraps() []error` case (non-standard interface) has no known caller in the stdlib or this codebase.
- `define_impl.go:Collect` — the `ErrorDefined`-passed-to-addf branch (90% covered) would require a separate test exercising that narrow sub-path.
