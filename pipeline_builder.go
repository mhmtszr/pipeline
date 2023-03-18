package pipeline

type Builder[K any] struct {
	steps []func(next StepDelegate[K]) StepDelegate[K]
}

func (t Builder[K]) build() Pipeline[K] {
	var step StepDelegate[K] = func(context K) {}
	for i := len(t.steps) - 1; i >= 0; i-- {
		step = t.steps[i](step)
	}
	return Pipeline[K]{
		stepDelegate: step,
	}
}

func (t Builder[K]) usePipelineStep(step Step[K]) Builder[K] {
	t.steps = append(t.steps, func(next StepDelegate[K]) StepDelegate[K] {
		return func(context K) {
			step.execute(context, next)
		}
	})
	return t
}
