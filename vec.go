package fn

import (
	"iter"
	"slices"
)

// Vec is a growable vector
type Vec[T any] []T

func Collect[T any](i iter.Seq[T]) Vec[T] {
	return slices.Collect(i)
}

// Iter implements [Iterable].
func (v Vec[T]) Iter() iter.Seq[T] {
	return slices.Values(v)
}

var _ Iterable[int] = make(Vec[int], 0)
