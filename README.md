# Pipeline [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][go-report-img]][go-report]

Go pipeline solution that can be used in many different combinations for chaining pipeline steps.

Inspired by [@bilal-kilic](https://github.com/bilal-kilic)'s Kotlin implementation [boru](https://github.com/Trendyol/boru).

### Usage
Supports 1.18+ Go versions because of Go Generics
```
go get github.com/mhmtszr/pipeline
```
### Examples

``` kotlin
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

[doc-img]: https://godoc.org/github.com/mhmtszr/pipeline?status.svg
[doc]: https://godoc.org/github.com/mhmtszr/pipeline
[ci-img]: https://github.com/mhmtszr/pipeline/actions/workflows/build-test.yml/badge.svg
[ci]: https://github.com/mhmtszr/pipeline/actions/workflows/build-test.yml
[cov-img]: https://codecov.io/gh/mhmtszr/pipeline/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/mhmtszr/pipeline
[go-report-img]: https://goreportcard.com/badge/github.com/mhmtszr/pipeline
[go-report]: https://goreportcard.com/report/github.com/mhmtszr/pipeline

