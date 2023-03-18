package pipeline

type Pipeline[K any] struct {
	stepDelegate StepDelegate[K]
}

func (t Pipeline[K]) Execute(context K) {
	t.stepDelegate(context)
}
