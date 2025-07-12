package errtrace

import "fmt"

type VarPoint struct {
	stacktrace []StackFrame
	vars map[string]any
}

type TracedError struct {
	stacktrace []StackFrame
	err        error
	vars       []VarPoint

	formatter func (stacktrace []StackFrame, err error, vars []VarPoint) string
}

func (e *TracedError) Error() string {
	return fmt.Sprintf("Error at %s:%d: %v", e.stacktrace[0].file, e.stacktrace[0].line, e.err)
}

func (e *TracedError) ErrorFormated() string {
	return e.formatter(e.stacktrace, e.err, e.vars)
}

func (e *TracedError) Unwrap() error {
	return e.err
}

