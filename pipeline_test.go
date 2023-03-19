package pipeline

import (
	"testing"
)

type Square struct{}
type Add struct{}

func (s Square) Execute(context *int, next func(context *int)) {
	*context = (*context) * (*context)
	next(context)
}

func (a Add) Execute(context *int, next func(context *int)) {
	*context = (*context) + (*context)
	next(context)
}

func TestPipeline(t *testing.T) {
	p := Builder[*int]{}.UsePipelineStep(Square{}).UsePipelineStep(Add{}).Build()
	nm := 3
	want := 18
	p.Execute(&nm)

	if nm != 18 {
		t.Errorf("got %d, wanted %d", nm, want)
	}
}
