// Package stack captures formatted call stacks.
package stack

import (
	"fmt"
	"runtime"
	"strings"
)

// maxFrames caps the number of program counters captured by Stack. A
// counter for inlined calls expands to one frame per inlined function, so
// Stack may return a few more than maxFrames frame strings.
const maxFrames = 32

// Stack returns the formatted frames of the calling goroutine's stack,
// innermost first, read from at most 32 program counters; inlined calls may
// expand them to a few more frame strings. skip is the number of frames to
// ascend, with 0 identifying the frame for Stack itself and 1 identifying the
// caller of Stack. Each frame is formatted as "file:line (0xpc): function".
func Stack(skip int) []string {
	var pcs [maxFrames]uintptr
	// runtime.Callers counts itself as frame 0, so skip+1 makes 0 mean Stack.
	n := runtime.Callers(skip+1, pcs[:])
	if n == 0 {
		return nil
	}

	s := make([]string, 0, n)
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		s = append(s, fmt.Sprintf("%s:%d (0x%x): %s", frame.File, frame.Line, frame.PC, function(frame.Function)))
		if !more {
			break
		}
	}

	return s
}

// function trims a fully qualified function name to its name within the
// package: "example.com/pkg.(*T).Method" becomes "(*T).Method".
func function(name string) string {
	if name == "" {
		return "???"
	}

	if slash := strings.LastIndex(name, "/"); slash >= 0 {
		name = name[slash+1:]
	}
	if dot := strings.Index(name, "."); dot >= 0 {
		name = name[dot+1:]
	}

	return strings.ReplaceAll(name, "·", ".")
}
