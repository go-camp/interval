## Interval

A golang package for performing operations on intervals.

## Usage

```go
package main

import (
	"fmt"

	"github.com/go-camp/interval"
)

func main() {
	a := &interval.OrderedSet{}
	a.Add(interval.Interval{
		Begin:    0,
		IncBegin: true,
		End:      10,
		IncEnd:   false,
	})
	a.Add(interval.Interval{
		Begin:    -10,
		IncBegin: true,
		End:      -5,
		IncEnd:   false,
	})
	fmt.Printf("a: %s\n", a)

	b := &interval.OrderedSet{}
	b.Add(interval.Interval{
		Begin:    -4,
		IncBegin: true,
		End:      5,
		IncEnd:   false,
	})
	fmt.Printf("b: %s\n", b)

	fmt.Printf("Union(a, b) = %s\n", interval.Union(a, b))
	fmt.Printf("Intersect(a, b) = %s\n", interval.Intersect(a, b))
	fmt.Printf("Subtract(a, b) = %s\n", interval.Subtract(a, b))
	fmt.Printf("Difference(a, b) = %s\n", interval.Difference(a, b))
	// a: {[-10, -5), [0, 10)}
	// b: {[-4, 5)}
	// Union(a, b) = {[-10, -5), [-4, 10)}
	// Intersect(a, b) = {[0, 5)}
	// Subtract(a, b) = {[-10, -5), [5, 10)}
	// Difference(a, b) = {[-10, -5), [-4, 0), [5, 10)}
}

```

## Install

```
go get -u -v github.com/go-camp/interval
```

## Test

```
go test ./...
```
