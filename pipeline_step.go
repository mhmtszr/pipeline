package pipeline

type PipelineStepDelegate[K any] func(context K)

type PipelineStep[K any] interface {
	execute(context K, next func(context K))
}
