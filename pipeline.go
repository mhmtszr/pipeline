package pipeline

type Pipeline[K any] struct {
	handler func(K) error
}

func (p Pipeline[K]) Execute(ctx K) error {
	return p.handler(ctx)
}
