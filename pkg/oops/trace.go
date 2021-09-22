package oops

import (
	"bytes"
	"fmt"
	"runtime"
)

var (
	traceUnknown = []byte("???")
	traceMidDot  = []byte("·")
	traceDot     = []byte(".")
	traceSlash   = []byte("/")
)

// stack returns stack information in a formatted string array
func stack(skip int) []string {
	s := make([]string, 0, 10)

	for idx := skip; ; idx++ {
		pc, file, line, ok := runtime.Caller(idx)
		if !ok {
			break
		}

		s = append(s,
			fmt.Sprintf("%s:%d (0x%x): %s", file, line, pc, function(pc)))
	}

	return s
}

func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return traceUnknown
	}
	name := []byte(fn.Name())

	if slash := bytes.LastIndex(name, traceSlash); slash >= 0 {
		name = name[slash+1:]
	}
	if dot := bytes.Index(name, traceDot); dot >= 0 {
		name = name[dot+1:]
	}

	return bytes.Replace(name, traceMidDot, traceDot, -1)
}
