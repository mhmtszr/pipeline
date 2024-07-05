package pipeline

type Builder[K any] struct {
	steps []func(next StepDelegate[K]) StepDelegate[K]
}

func (t Builder[K]) Build() Pipeline[K] {
	var step StepDelegate[K] = func(context K) error {
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
