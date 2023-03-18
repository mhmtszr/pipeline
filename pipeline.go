package pipeline

type Pipeline[K any] struct {
	pipelineStepDelegate PipelineStepDelegate[K]
}

func (t Pipeline[K]) execute(context K) {
	t.pipelineStepDelegate(context)
}
