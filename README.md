# Pipeline

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

