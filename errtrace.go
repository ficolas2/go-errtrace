package errtrace

import (
	"errors"
	"fmt"
)

type VarPoint struct {
	stacktrace []StackFrame
	vars map[string]any
}

type TracedError struct {
	stacktrace []StackFrame
	err        error
	vars       []VarPoint
}

func (e *TracedError) Error() string {
	return fmt.Sprintf("Error at %s:%d: %v", e.stacktrace[0].file, e.stacktrace[0].line, e.err)
}

func (e *TracedError) ErrorFormated() string {
	// TODO
	return ""
}

func (e *TracedError) Unwrap() error {
	return e.err
}

var defaultTracer Tracer = NewDefaultTracer()

func Wrap(err error) error {
	return defaultTracer.Wrap(err)
}

func WrapVars(err error, vars map[string]any) error {
	return defaultTracer.WrapVars(err, vars)
}

func (t *tracer) Wrap(err error) error {
	return t.WrapVars(err, nil)
}

func (t *tracer) WrapVars(err error, vars map[string]any) error {
	if err == nil {
		return nil
	}

	var tracedErr *TracedError
	if errors.As(err, &tracedErr) {
		stacktrace := captureStacktrace(2, t.maxVarStackDepth)
		appendArgs(tracedErr, vars, stacktrace[:t.maxVarStackDepth])
		return tracedErr
	}

	stacktrace := captureStacktrace(2, t.maxStackDepth)

	tracedErr = &TracedError {
		stacktrace: stacktrace,
		err:        err,
	}
	appendArgs(tracedErr, vars, stacktrace[:t.maxVarStackDepth])

	return tracedErr
}

func appendArgs(err *TracedError, vars map[string]any, stacktrace []StackFrame) {
	if len(vars) == 0 {
		return
	}

	err.vars = append(err.vars, VarPoint{
		stacktrace: stacktrace,
		vars: vars,
	})
}

