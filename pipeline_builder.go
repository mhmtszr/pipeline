package pipeline

type PipelineBuilder[K any] struct {
	steps []func(next PipelineStepDelegate[K]) PipelineStepDelegate[K]
}

func (t PipelineBuilder[K]) build() Pipeline[K] {
	var step PipelineStepDelegate[K] = func(context K) {}
	for i := len(t.steps) - 1; i >= 0; i-- {
		step = t.steps[i](step)
	}
	return Pipeline[K]{
		pipelineStepDelegate: step,
	}
}

func (t PipelineBuilder[K]) usePipelineStep(step PipelineStep[K]) PipelineBuilder[K] {
	t.steps = append(t.steps, func(next PipelineStepDelegate[K]) PipelineStepDelegate[K] {
		return func(context K) {
			step.execute(context, next)
		}
	})
	return t
}
