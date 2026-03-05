# Changelog


## v1.0.1 Released (2026-03-05)

* Add CLAUDE
* Bug fix — `oops.As` dead code path (`Unwarp` typo)
* Add tests for coverage

## v1.0.0 Released (2026-01-29)

* Includes `oops` as is, releasing the first `v1` tag.
* Migration from `v0` depends on how old the `v0.x` is:
  * `v0.12+` should be enough to just go update (point to v1.0.0)
  * `v0.10+` might require a little refactoring, but overall possible
  * `pre-v0.10` requires major refactoring
  * `pre-v0.3` complete refactoring
