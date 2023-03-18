package pipeline

type StepDelegate[K any] func(context K)

type Step[K any] interface {
	Execute(context K, next func(context K))
}
