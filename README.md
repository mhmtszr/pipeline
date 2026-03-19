# Pipeline [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov] [![Go Report Card][go-report-img]][go-report]

![133256800-8d51f5e5-1cc5-45d2-a195-28e95f1cb92c](https://user-images.githubusercontent.com/35079195/226874057-bc579acd-1f0a-4dc7-85dc-8ec73811e56e.jpeg)


Go pipeline solution that can be used in many different combinations for chaining pipeline steps.

Inspired by [@bilal-kilic](https://github.com/bilal-kilic)'s Kotlin implementation [boru](https://github.com/Trendyol/boru).

### Usage
Supports 1.22+ Go versions because of Go Generics
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

func square(ctx *int, next func(*int) error) error {
	*ctx = (*ctx) * (*ctx)
	fmt.Printf("After first step: %d\n", *ctx)
	return next(ctx)
}

func add(ctx *int, next func(*int) error) error {
	*ctx = (*ctx) + (*ctx)
	fmt.Printf("After second step: %d\n", *ctx)
	return next(ctx)
}

func main() {
	p := pipeline.NewBuilder[*int]().Use(square).Use(add).Build()
	nm := 3
	_ = p.Execute(&nm)
}
// After first step: 9
// After second step: 18
```

#### Concurrent Pipeline

``` go
p := pipeline.NewBuilder[*atomic.Uint64]().
	UseConcurrent(
		func(ctx *atomic.Uint64) error {
			ctx.Add(ctx.Load())
			return nil
		},
		func(ctx *atomic.Uint64) error {
			ctx.Add(ctx.Load())
			return nil
		},
	).Build()

var nmb atomic.Uint64
nmb.Add(5)
_ = p.Execute(&nmb)
```

#### Conditional Pipeline

``` go
p := pipeline.NewBuilder[*int]().
	UseConditional(
		func(ctx *int) bool { return *ctx == 3 },
		[]pipeline.StepFunc[*int]{square},
		[]pipeline.StepFunc[*int]{add},
	).Use(add).Build()

nm := 3
_ = p.Execute(&nm)
// nm == 18

nm = 4
_ = p.Execute(&nm)
// nm == 16
```
[doc-img]: https://godoc.org/github.com/mhmtszr/pipeline?status.svg
[doc]: https://godoc.org/github.com/mhmtszr/pipeline
[ci-img]: https://github.com/mhmtszr/pipeline/actions/workflows/build-test.yml/badge.svg
[ci]: https://github.com/mhmtszr/pipeline/actions/workflows/build-test.yml
[cov-img]: https://codecov.io/gh/mhmtszr/pipeline/branch/master/graph/badge.svg
[cov]: https://codecov.io/gh/mhmtszr/pipeline
[go-report-img]: https://goreportcard.com/badge/github.com/mhmtszr/pipeline
[go-report]: https://goreportcard.com/report/github.com/mhmtszr/pipeline
