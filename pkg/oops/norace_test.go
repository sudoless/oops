//go:build !race

package oops_test

// raceEnabled reports whether the race detector is on; it changes allocation counts.
const raceEnabled = false
