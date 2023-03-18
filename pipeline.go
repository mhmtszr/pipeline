package pipeline

type Pipeline[K any] struct {
	stepDelegate StepDelegate[K]
}

func (t Pipeline[K]) execute(context K) {
	t.stepDelegate(context)
}
