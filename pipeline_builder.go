package pipeline

import "golang.org/x/sync/errgroup"

type Builder[K any] struct {
	steps []StepFunc[K]
}

func NewBuilder[K any]() *Builder[K] {
	return &Builder[K]{}
}

func (b *Builder[K]) Use(step StepFunc[K]) *Builder[K] {
	b.steps = append(b.steps, step)
	return b
}

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

func (b *Builder[K]) Build() Pipeline[K] {
	return Pipeline[K]{handler: buildChain(b.steps)}
}

func buildChain[K any](steps []StepFunc[K]) func(K) error {
	var h = func(K) error { return nil }
	for i := len(steps) - 1; i >= 0; i-- {
		next := h
		step := steps[i]
		h = func(ctx K) error { return step(ctx, next) }
	}
	return h
}
