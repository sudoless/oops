package unsafe

import (
	"bytes"
	"fmt"
	"runtime"
)

var (
	stackUnknown = []byte("???")
	stackMidDot  = []byte("·")
	stackDot     = []byte(".")
	stackSlash   = []byte("/")
)

// Stack returns stack information in a formatted string slice.
func Stack(skip int) []string {
	s := make([]string, 0, 10)

	for idx := skip; ; idx++ {
		pc, file, line, ok := runtime.Caller(idx)
		if !ok {
			break
		}

		s = append(s,
			fmt.Sprintf("%s:%d (0x%x): %s", file, line, pc, stackFunction(pc)))
	}

	return s
}

func stackFunction(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return stackUnknown
	}
	name := []byte(fn.Name())

	if slash := bytes.LastIndex(name, stackSlash); slash >= 0 {
		name = name[slash+1:]
	}
	if dot := bytes.Index(name, stackDot); dot >= 0 {
		name = name[dot+1:]
	}

	return bytes.Replace(name, stackMidDot, stackDot, -1)
}
