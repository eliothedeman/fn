# fn

A Go package providing functional programming utilities built on Go's `iter.Seq` iterators. Uses a lisp-style free function API.

[![Go](https://github.com/eliothedeman/fn/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/eliothedeman/fn/actions/workflows/go.yml)
[![GoDoc](https://pkg.go.dev/badge/github.com/eliothedeman/fn.svg)](https://pkg.go.dev/github.com/eliothedeman/fn)

## Requirements

- Go 1.24+

## Install

```
go get github.com/eliothedeman/fn
```

## Features

### Iterators

- **`Range(start, end)`** — produces a sequence of values from `start` to `end` (exclusive) with a step of 1
- **`StepRange(start, end, step)`** — produces a sequence of values from `start` to `end` (exclusive) with a custom step
- **`Chain(iters...)`** — concatenates multiple iterators into a single sequence
- **`Map(iter, f)`** — transforms each element in an iterator using a function
- **`Filter(iter, pred)`** — yields only elements that satisfy a predicate
- **`Reduce(iter, seed, f)`** — folds an iterator into a single value
- **`Sum(iter)`** — sums all numeric values in an iterator

### Shared interfaces

Both `Result[T]` and `Option[T]` satisfy `Iterable[T]` and work with the same set of unwrap functions:

- **`Iter(x)`** — returns the iterator from any `Iterable[T]`
- **`HasValue(x)`** — returns true if the container holds a value
- **`IsEmpty(x)`** — returns true if the container is empty (Err or None)
- **`Unwrap(x)`** — returns the value or panics if empty/error
- **`UnwrapOr(x, def)`** — returns the value or a default
- **`UnwrapOrF(x, f)`** — returns the value or calls a function to produce a default

### Result

A generic `Result[T]` type for representing a value-or-error.

- **`Ok(val)`** — creates a successful result
- **`Err[T](err)`** — creates an error result
- **`Try(val, err)`** — creates a result from a `(T, error)` pair, common with Go APIs
- **`Unpack(r)`** — returns the `(T, error)` pair
- **`IterErr(r)`** — yields the error if Err, nothing if Ok

### Option

A generic `Option[T]` type for representing an optional value.

- **`Some(val)`** — creates an Option containing a value
- **`None[T]()`** — creates an empty Option

## Usage

```go
package main

import (
	"fmt"

	"github.com/eliothedeman/fn"
)

func main() {
	// Sum integers 0..99
	total := fn.Sum(fn.Range(0, 100))
	fmt.Println(total) // 4950

	// Map and filter
	evens := fn.Filter(fn.Range(0, 20), func(i int) bool {
		return i%2 == 0
	})
	doubled := fn.Map(evens, func(i int) int {
		return i * 2
	})
	fmt.Println(fn.Sum(doubled)) // 180

	// Result type
	r := fn.Ok(42)
	val, err := fn.Unpack(r)
	fmt.Println(val, err) // 42 <nil>

	// Option type
	o := fn.Some("hello")
	fmt.Println(fn.UnwrapOr(o, "default")) // hello
	fmt.Println(fn.UnwrapOr(fn.None[string](), "default")) // default
}
```

## License

MIT
