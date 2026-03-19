package pipeline_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/mhmtszr/pipeline"
)

func square(ctx *int, next func(*int) error) error {
	*ctx = (*ctx) * (*ctx)
	return next(ctx)
}

func add(ctx *int, next func(*int) error) error {
	*ctx = (*ctx) + (*ctx)
	return next(ctx)
}

func multiply(_ *int, _ func(*int) error) error {
	return fmt.Errorf("multiplyStepError")
}

func addConcurrent(ctx *atomic.Uint64) error {
	ctx.Add(ctx.Load())
	return nil
}

func squareConcurrentErr(_ *atomic.Uint64) error {
	return fmt.Errorf("errorTest")
}

func TestPipeline(t *testing.T) {
	p := pipeline.NewBuilder[*int]().Use(square).Use(add).Build()
	nm := 3
	want := 18
	err := p.Execute(&nm)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if nm != want {
		t.Errorf("got %d, wanted %d", nm, want)
	}
}

func TestConditionalPipeline(t *testing.T) {
	p := pipeline.NewBuilder[*int]().
		UseConditional(
			func(ctx *int) bool { return *ctx == 3 },
			[]pipeline.StepFunc[*int]{square},
			[]pipeline.StepFunc[*int]{add},
		).Use(add).Build()

	nm := 3
	want := 18
	_ = p.Execute(&nm)
	if nm != want {
		t.Errorf("got %d, wanted %d", nm, want)
	}

	nm = 4
	want = 16
	_ = p.Execute(&nm)
	if nm != want {
		t.Errorf("got %d, wanted %d", nm, want)
	}
}

func TestConditionalPipelineStepError(t *testing.T) {
	p := pipeline.NewBuilder[*int]().
		UseConditional(
			func(ctx *int) bool { return *ctx == 3 },
			[]pipeline.StepFunc[*int]{multiply},
			nil,
		).Use(add).Build()

	nm := 3
	wantErr := "multiplyStepError"
	err := p.Execute(&nm)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != wantErr {
		t.Errorf("got %s, wanted %s", err.Error(), wantErr)
	}
}

func TestConcurrentStep(t *testing.T) {
	p := pipeline.NewBuilder[*atomic.Uint64]().
		UseConcurrent(addConcurrent, addConcurrent).
		Build()

	want := uint64(20)
	var nmb atomic.Uint64
	nmb.Add(5)
	_ = p.Execute(&nmb)

	if load := nmb.Load(); load != want {
		t.Errorf("got %d, wanted %d", load, want)
	}
}

func TestConcurrentStepError(t *testing.T) {
	p := pipeline.NewBuilder[*atomic.Uint64]().
		UseConcurrent(addConcurrent, squareConcurrentErr).
		Build()

	var nmb atomic.Uint64
	nmb.Add(5)
	err := p.Execute(&nmb)

	wantErr := "errorTest"
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != wantErr {
		t.Errorf("got %s, wanted %s", err.Error(), wantErr)
	}
}
