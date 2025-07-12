package errtrace

import (
	"errors"
)

var defaultTracer Tracer = NewDefaultTracer()

// Entry points to the wrapping. The skip is very important, so that no method from this
// library gets added to the stack
func SetDefaultTracer(tracer Tracer) {
	defaultTracer = tracer
}

func Wrap(err error) error {
	return defaultTracer.wrapInternal(err, 1)
}

func WrapVars(err error, vars map[string]any) error {
	return defaultTracer.wrapVarsInternal(err, vars, 1)
}

func (t *tracer) Wrap(err error) error {
	return t.wrapVarsInternal(err, nil, 1)
}

func (t *tracer) WrapVars(err error, vars map[string]any) error {
	return t.wrapVarsInternal(err, vars, 1)
}

// Internal wrap functions
func (t *tracer) wrapInternal(err error, skip int) error {
	return t.wrapVarsInternal(err, nil, skip + 1)
}

func (t *tracer) wrapVarsInternal(err error, vars map[string]any, skip int) error {
	if err == nil {
		return nil
	}

	var tracedErr *TracedError
	if errors.As(err, &tracedErr) {
		stacktrace := captureStacktrace(skip + 1, t.maxVarStackDepth, t.trimStackTrace)
		appendArgs(tracedErr, vars, stacktrace[:t.maxVarStackDepth])
		return tracedErr
	}

	stacktrace := captureStacktrace(skip + 1, t.maxStackDepth, t.trimStackTrace)

	tracedErr = &TracedError {
		stacktrace: stacktrace,
		err:        err,
		formatter: t.formatter,
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

