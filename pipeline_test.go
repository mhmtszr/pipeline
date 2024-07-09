package pipeline_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/mhmtszr/pipeline"
)

type (
	Square   struct{}
	Add      struct{}
	Multiply struct{}
)

func (s Square) Execute(context *int, next func(context *int) error) error {
	*context = (*context) * (*context)
	return next(context)
}

func (m Multiply) Execute(_ *int, _ func(context *int) error) error {
	return fmt.Errorf("multiplyStepError")
}

func (s Square) ConcurrentExecute(_ *atomic.Uint64) error {
	return fmt.Errorf("errorTest")
}

func (a Add) Execute(context *int, next func(context *int) error) error {
	*context = (*context) + (*context)
	return next(context)
}

func (a Add) ConcurrentExecute(context *atomic.Uint64) error {
	context.Add(context.Load())
	return nil
}

func TestPipeline(t *testing.T) {
	p := pipeline.Builder[*int]{}.UsePipelineStep(Square{}).UsePipelineStep(Add{}).Build()
	nm := 3
	want := 18
	err := p.Execute(&nm)
	if err != nil {
		t.Errorf("got error %s", err.Error())
	}

	if nm != 18 {
		t.Errorf("got %d, wanted %d", nm, want)
	}
}

func TestConditionalPipeline(t *testing.T) {
	p := pipeline.Builder[*int]{}.
		UseConditionalStepBuilder(
			pipeline.NewConditionalStepBuilder[*int]().
				Condition(func(context *int) bool {
					return *context == 3
				}).
				IfTrue(Square{}).
				IfFalse(Add{}),
		).UsePipelineStep(Add{}).Build()
	nm := 3
	want := 18
	_ = p.Execute(&nm)

	if nm != 18 {
		t.Errorf("got %d, wanted %d", nm, want)
	}

	nm = 4
	want = 16
	_ = p.Execute(&nm)

	if nm != 16 {
		t.Errorf("got %d, wanted %d", nm, want)
	}
}

func TestConditionalPipelineStepError(t *testing.T) {
	p := pipeline.Builder[*int]{}.
		UseConditionalStepBuilder(
			pipeline.NewConditionalStepBuilder[*int]().
				Condition(func(context *int) bool {
					return *context == 3
				}).
				IfTrue(Multiply{}),
		).UsePipelineStep(Add{}).Build()
	nm := 3
	wantErr := "multiplyStepError"
	err := p.Execute(&nm)

	if err == nil || err.Error() != wantErr {
		t.Errorf("got %s, wanted %s", err.Error(), wantErr)
	}
}

func TestConditionalPipelineErrorWithoutCondition(t *testing.T) {
	p := pipeline.Builder[*int]{}.
		UseConditionalStepBuilder(
			pipeline.NewConditionalStepBuilder[*int]().
				IfTrue(Square{}).
				IfFalse(Add{}),
		).UsePipelineStep(Add{}).Build()
	nm := 3
	err := p.Execute(&nm)

	wantErr := "condition not found in conditional step"

	if err == nil || err.Error() != wantErr {
		t.Errorf("got %s, wanted %s", err.Error(), wantErr)
	}
}

func TestConcurrentStep(t *testing.T) {
	p := pipeline.Builder[*atomic.Uint64]{}.
		UseConcurrentPipelineSteps(
			Add{}, Add{},
		).Build()

	want := 20
	var nmb atomic.Uint64
	nmb.Add(5)
	_ = p.Execute(&nmb)

	load := nmb.Load()
	if load != uint64(want) {
		t.Errorf("got %d, wanted %d", load, want)
	}
}

func TestConcurrentStepError(t *testing.T) {
	p := pipeline.Builder[*atomic.Uint64]{}.
		UseConcurrentPipelineSteps(
			Add{}, Square{},
		).Build()

	var nmb atomic.Uint64
	nmb.Add(5)
	err := p.Execute(&nmb)

	wantErr := "errorTest"
	if err == nil || err.Error() != wantErr {
		t.Errorf("got %s, wanted %s", err.Error(), wantErr)
	}
}
