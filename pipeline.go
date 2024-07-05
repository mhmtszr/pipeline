package pipeline

type Pipeline[K any] struct {
	stepDelegate StepDelegate[K]
}

func (t Pipeline[K]) Execute(context K) error {
	return t.stepDelegate(context)
}
