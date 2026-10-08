# CLAUDE.md

1. The repo is public. Commit messages, PRs and docs describe API and behaviour, never decision history.
2. The core module has zero dependencies.
3. Nothing adds work to the `Yeet`, `Wrap` or `Error()` hot paths. Diagnostics belong in tests or tooling.
4. Every `*Error` method is nil-receiver safe, and every intake treats a typed-nil `*Error` as nil.
5. Package helpers that pass an error along return `error`, never `*Error`.
