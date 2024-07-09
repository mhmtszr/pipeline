# Pipeline [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][go-report-img]][go-report]

![133256800-8d51f5e5-1cc5-45d2-a195-28e95f1cb92c](https://user-images.githubusercontent.com/35079195/226874057-bc579acd-1f0a-4dc7-85dc-8ec73811e56e.jpeg)


Go pipeline solution that can be used in many different combinations for chaining pipeline steps.

Inspired by [@bilal-kilic](https://github.com/bilal-kilic)'s Kotlin implementation [boru](https://github.com/Trendyol/boru).

### Usage
Supports 1.18+ Go versions because of Go Generics
```
go get github.com/mhmtszr/pipeline
```
### Examples
#### Basic Pipeline
``` go
package main

import (
	"fmt"
	"github.com/mhmtszr/pipeline"
)

type Square struct{}
type Add struct{}

func (s Square) Execute(context int, next func(context int)) {
	context = context * context
	println(fmt.Sprintf("After first chain context: %d", context))
	next(context)
}

func (a Add) Execute(context int, next func(context int)) {
	context = context + context
	println(fmt.Sprintf("After second chain context: %d", context))
	next(context)
}

func main() {
	p := pipeline.Builder[int]{}.UsePipelineStep(Square{}).UsePipelineStep(Add{}).Build()
	p.Execute(3)
}
// After first chain context: 9
// After second chain context: 18

```

#### Conditional Pipeline

``` go
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
	_ = p.Execute(&nm)

	// nm 18

	nm = 4
	_ = p.Execute(&nm)

	// nm 16
```
[doc-img]: https://godoc.org/github.com/mhmtszr/pipeline?status.svg
[doc]: https://godoc.org/github.com/mhmtszr/pipeline
[ci-img]: https://github.com/mhmtszr/pipeline/actions/workflows/build-test.yml/badge.svg
[ci]: https://github.com/mhmtszr/pipeline/actions/workflows/build-test.yml
[cov-img]: https://codecov.io/gh/mhmtszr/pipeline/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/mhmtszr/pipeline
[go-report-img]: https://goreportcard.com/badge/github.com/mhmtszr/pipeline
[go-report]: https://goreportcard.com/report/github.com/mhmtszr/pipeline

