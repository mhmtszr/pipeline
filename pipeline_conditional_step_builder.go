package pipeline

type ConditionalStepBuilder[K any] struct {
	condition  func(context K) bool
	trueSteps  []Step[K]
	falseSteps []Step[K]
}

func (b *ConditionalStepBuilder[K]) Condition(condition func(context K) bool) *ConditionalStepBuilder[K] {
	b.condition = condition
	return b
}

func (b *ConditionalStepBuilder[K]) IfTrue(steps ...Step[K]) *ConditionalStepBuilder[K] {
	b.trueSteps = steps
	return b
}

func (b *ConditionalStepBuilder[K]) IfFalse(steps ...Step[K]) *ConditionalStepBuilder[K] {
	b.falseSteps = steps
	return b
}

func (b *ConditionalStepBuilder[K]) Build() func(next StepDelegate[K]) StepDelegate[K] {
	return func(next StepDelegate[K]) StepDelegate[K] {
		truePipeline := Builder[K]{steps: b.toStepDelegates(b.trueSteps)}.Build()
		falsePipeline := Builder[K]{steps: b.toStepDelegates(b.falseSteps)}.Build()
		return func(context K) error {
			return ConditionalStep[K]{b.condition, truePipeline, falsePipeline}.Execute(context, next)
		}
	}
}

func (b *ConditionalStepBuilder[K]) toStepDelegates(steps []Step[K]) []func(next StepDelegate[K]) StepDelegate[K] {
	stepDelegates := make([]func(next StepDelegate[K]) StepDelegate[K], len(steps))
	for i, step := range steps {
		stepDelegates[i] = func(next StepDelegate[K]) StepDelegate[K] {
			return func(context K) error {
				return step.Execute(context, next)
			}
		}
	}
	return stepDelegates
}

func NewConditionalStepBuilder[K any]() *ConditionalStepBuilder[K] {
	return &ConditionalStepBuilder[K]{}
}
