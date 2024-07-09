package pipeline

import "golang.org/x/sync/errgroup"

type Builder[K any] struct {
	steps []func(next StepDelegate[K]) StepDelegate[K]
}

func (t Builder[K]) Build() Pipeline[K] {
	var step StepDelegate[K] = func(_ K) error {
		return nil
	}
	for i := len(t.steps) - 1; i >= 0; i-- {
		step = t.steps[i](step)
	}
	return Pipeline[K]{
		stepDelegate: step,
	}
}

func (t Builder[K]) UsePipelineStep(step Step[K]) Builder[K] {
	t.steps = append(t.steps, func(next StepDelegate[K]) StepDelegate[K] {
		return func(context K) error {
			return step.Execute(context, next)
		}
	})
	return t
}

func (t Builder[K]) UseConcurrentPipelineSteps(steps ...ConcurrentStep[K]) Builder[K] {
	t.steps = append(t.steps, func(next StepDelegate[K]) StepDelegate[K] {
		return func(context K) error {
			var eg errgroup.Group
			for _, step := range steps {
				step := step
				eg.Go(func() error {
					return step.ConcurrentExecute(context)
				})
			}
			err := eg.Wait()
			if err != nil {
				return err
			}
			return next(context)
		}
	})
	return t
}

func (t Builder[K]) UseConditionalStepBuilder(builder *ConditionalStepBuilder[K]) Builder[K] {
	t.steps = append(t.steps, builder.Build())
	return t
}
