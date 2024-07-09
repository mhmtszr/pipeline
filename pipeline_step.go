package pipeline

type StepDelegate[K any] func(context K) error

type Step[K any] interface {
	Execute(context K, next func(context K) error) error
}

type ConcurrentStep[K any] interface {
	ConcurrentExecute(context K) error
}
