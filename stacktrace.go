package errtrace

import (
	"runtime"
	"strings"
)

type StackFrame struct {
	file string
	line int
	function string
}

func captureStacktrace(skip int, maxFrames int, trimStacktrace *string) []StackFrame {
	pcs := make([]uintptr, maxFrames)
	n := runtime.Callers(skip+2, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	var stack []StackFrame
	for {
		frame, more := frames.Next()
		stack = append(stack, StackFrame{
			file:     frame.File,
			line:     frame.Line,
			function: frame.Function,
		})
		if strings.Contains(frame.File, *trimStacktrace) {
			break
		}
		if !more {
			break
		}
	}

	return stack
}
