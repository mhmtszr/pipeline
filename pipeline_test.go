package pipeline_test

import (
	"github.com/mhmtszr/pipeline"
	"testing"
)

type (
	Square struct{}
	Add    struct{}
)

func (s Square) Execute(context *int, next func(context *int) error) error {
	*context = (*context) * (*context)
	return next(context)
}

func (a Add) Execute(context *int, next func(context *int) error) error {
	*context = (*context) + (*context)
	return next(context)
}

func TestPipeline(t *testing.T) {
	p := pipeline.Builder[*int]{}.UsePipelineStep(Square{}).UsePipelineStep(Add{}).Build()
	nm := 3
	want := 18
	p.Execute(&nm)

	if nm != 18 {
		t.Errorf("got %d, wanted %d", nm, want)
	}
}
