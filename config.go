package errtrace

import "errors"

type Tracer interface {
	Wrap(err error) error
	WrapVars(err error, vars map[string]any) error
}

type tracer struct {
	maxStackDepth    int
	maxVarStackDepth int

	formatter func(stacktrace []StackFrame, err error, vars []VarPoint) string
}

type TracerBuilder struct {
	tracer tracer
}

func NewDefaultTracer() Tracer {
	return NewTracerBuilder().MustBuild()
}

func NewTracerBuilder() *TracerBuilder {
	return &TracerBuilder{
		tracer: tracer{
			maxStackDepth:    32,
			maxVarStackDepth: 1,
			formatter:        defaultFormatErr,
		},
	}
}

func (t *TracerBuilder) MustBuild() Tracer {
	tracer, err := t.Build()
	if err != nil {
		panic(err)
	}

	return tracer
}

func (t *TracerBuilder) Build() (Tracer, error) {
	if t.tracer.maxVarStackDepth > t.tracer.maxStackDepth {
		return nil, errors.New("MaxVarStackDepth cannot be bigger than MaxStackDepth")
	}
	if t.tracer.maxVarStackDepth < 1 {
		return nil, errors.New("MaxVarStackDepth cannot be less than one")
	}
	if t.tracer.maxStackDepth < 1 {
		return nil, errors.New("MaxStackDepth cannot be less than one")
	}

	return &t.tracer, nil
}

func (t *TracerBuilder) MaxVarStackDepth(m int) *TracerBuilder {
	t.tracer.maxStackDepth = m
	return t
}

func (t *TracerBuilder) MaxStackDepth(m int) *TracerBuilder {
	t.tracer.maxVarStackDepth = m
	return t
}

func (t *TracerBuilder) SetFormatter(
	formatter func(stacktrace []StackFrame, err error, vars []VarPoint) string,
) *TracerBuilder {
	t.tracer.formatter = formatter
	return t
}
