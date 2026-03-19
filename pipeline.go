// Package pipeline provides a generic, function-based pipeline builder.
package pipeline

// Pipeline is a compiled chain of steps ready to execute.
type Pipeline[K any] struct {
	handler func(K) error
}

// Execute runs the pipeline with the given context.
func (p Pipeline[K]) Execute(ctx K) error {
	return p.handler(ctx)
}
