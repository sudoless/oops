package oops

import (
	"fmt"
	"runtime"
	"strings"
)

// traceFrameName strips the module path prefix from a fully-qualified function
// name, returning "package.Function" or just "Function".
func traceFrameName(name string) string {
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.Index(name, "."); idx >= 0 {
		name = name[idx+1:]
	}

	return strings.ReplaceAll(name, "·", ".")
}

// Trace returns the formatted stack frames captured at creation time.
// Frames are formatted lazily on each call; returns nil if no trace was captured.
func (err *Error) Trace() []string {
	if len(err.trace) == 0 {
		return nil
	}

	frames := runtime.CallersFrames(err.trace)
	result := make([]string, 0, len(err.trace))
	for {
		frame, more := frames.Next()
		result = append(result, fmt.Sprintf("%s:%d (0x%x): %s",
			frame.File, frame.Line, frame.PC, traceFrameName(frame.Function)))
		if !more {
			break
		}
	}
	return result
}

// TraceFrames returns the raw stack frame pointers captured at creation time.
func (err *Error) TraceFrames() []uintptr {
	return err.trace
}
