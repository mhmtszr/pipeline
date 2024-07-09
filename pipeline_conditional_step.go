package pipeline

import "errors"

type ConditionalStep[K any] struct {
	condition     func(context K) bool
	truePipeline  Pipeline[K]
	falsePipeline Pipeline[K]
}

var errConditionNotFound = errors.New("condition not found in conditional step")

func (c ConditionalStep[K]) Execute(context K, next func(context K) error) error {
	if c.condition == nil {
		return errConditionNotFound
	}

	var err error
	if c.condition(context) {
		err = c.truePipeline.Execute(context)
	} else {
		err = c.falsePipeline.Execute(context)
	}
	if err != nil {
		return err
	}
	return next(context)
}
