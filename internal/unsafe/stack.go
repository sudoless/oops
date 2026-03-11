package stack

import "runtime"

const maxDepth = 64

// Stack returns raw program counter values for the call stack, starting skip
// frames above the caller. Use runtime.CallersFrames to format the result.
func Stack(skip int) []uintptr {
	pcs := make([]uintptr, maxDepth)
	// skip+1: runtime.Callers itself is frame 0; skip the Stack frame plus
	// the requested number of caller frames above it.
	n := runtime.Callers(skip+1, pcs)
	return pcs[:n]
}
