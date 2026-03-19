package pipeline

// StepFunc represents a pipeline step that receives context and a next function.
type StepFunc[K any] func(ctx K, next func(K) error) error
