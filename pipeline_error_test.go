package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/mhmtszr/pipeline"
)

type (
	SuccessStep struct{}
	ErrorStep   struct{}
)

func (s SuccessStep) Execute(context *int, next func(context *int) error) error {
	return next(context)
}

func (e ErrorStep) Execute(context *int, next func(context *int) error) error {
	return fmt.Errorf("errorstep error")
}

func TestErrorPipeline(t *testing.T) {
	p := pipeline.Builder[*int]{}.UsePipelineStep(SuccessStep{}).UsePipelineStep(ErrorStep{}).Build()
	nm := 3
	err := p.Execute(&nm)
	if err != nil && err.Error() == "errorstep error" {
		return
	}
	t.Errorf("error step should return error")
}
