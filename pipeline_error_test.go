package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/mhmtszr/pipeline"
)

func successStep(ctx *int, next func(*int) error) error {
	return next(ctx)
}

func errorStep(_ *int, _ func(*int) error) error {
	return fmt.Errorf("errorstep error")
}

func TestErrorPipeline(t *testing.T) {
	p := pipeline.NewBuilder[*int]().Use(successStep).Use(errorStep).Build()
	nm := 3
	wantErr := "errorstep error"
	err := p.Execute(&nm)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != wantErr {
		t.Errorf("got %s, wanted %s", err.Error(), wantErr)
	}
}
