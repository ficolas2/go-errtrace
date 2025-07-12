package errtrace

import "errors"

type Tracer interface {
	Wrap(err error) error
	WrapVars(err error, vars map[string]any) error
}

type tracer struct {
	maxStackDepth int
	maxVarStackDepth int
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

func (t *TracerBuilder) MaxVarStackDepth(m int) {
	t.tracer.maxStackDepth = m
}

func (t *TracerBuilder) MaxStackDepth(m int) {
	t.tracer.maxVarStackDepth = m
}
