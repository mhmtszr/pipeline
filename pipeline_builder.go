package pipeline

import "golang.org/x/sync/errgroup"

// Builder constructs a pipeline by chaining steps.
type Builder[K any] struct {
	steps []StepFunc[K]
}

// NewBuilder creates a new pipeline builder.
func NewBuilder[K any]() *Builder[K] {
	return &Builder[K]{}
}

// Use adds a step to the pipeline.
func (b *Builder[K]) Use(step StepFunc[K]) *Builder[K] {
	b.steps = append(b.steps, step)
	return b
}

// UseConcurrent adds a set of functions that run concurrently as a single step.
func (b *Builder[K]) UseConcurrent(steps ...func(K) error) *Builder[K] {
	b.steps = append(b.steps, func(ctx K, next func(K) error) error {
		var eg errgroup.Group
		for _, s := range steps {
			eg.Go(func() error { return s(ctx) })
		}
		if err := eg.Wait(); err != nil {
			return err
		}
		return next(ctx)
	})
	return b
}

// UseConditional adds a conditional step that branches based on cond.
func (b *Builder[K]) UseConditional(
	cond func(K) bool,
	ifTrue []StepFunc[K],
	ifFalse []StepFunc[K],
) *Builder[K] {
	trueH := buildChain(ifTrue)
	falseH := buildChain(ifFalse)
	b.steps = append(b.steps, func(ctx K, next func(K) error) error {
		h := falseH
		if cond(ctx) {
			h = trueH
		}
		if err := h(ctx); err != nil {
			return err
		}
		return next(ctx)
	})
	return b
}

// Build compiles the steps into an executable Pipeline.
func (b *Builder[K]) Build() Pipeline[K] {
	return Pipeline[K]{handler: buildChain(b.steps)}
}

func buildChain[K any](steps []StepFunc[K]) func(K) error {
	h := func(K) error { return nil }
	for i := len(steps) - 1; i >= 0; i-- {
		next := h
		step := steps[i]
		h = func(ctx K) error { return step(ctx, next) }
	}
	return h
}
