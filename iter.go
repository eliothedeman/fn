package fn

import (
	"iter"

	"golang.org/x/exp/constraints"
)

type Result[T any] struct {
	val T
	err error
}

func Try[T any](t T, err error) Result[T] {
	return Result[T]{val: t, err: err}
}

func Ok[T any](t T) Result[T] {
	return Result[T]{val: t}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

func (r Result[T]) Unpack() (T, error) {
	return r.val, r.err
}

func (r Result[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if r.err == nil {
			yield(r.val)
		}
	}
}

func (r Result[T]) IsOk() bool {
	return r.err == nil
}

func (r Result[T]) IsErr() bool {
	return r.err != nil
}

func (r *Result[T]) IterErr() iter.Seq[error] {
	return func(yield func(error) bool) {
		if r.err != nil {
			yield(r.err)
		}
	}
}

type Option[T any] struct {
	val     T
	hasSome bool
}

func (o *Option[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.hasSome {
			yield(o.val)
		}
	}
}

func (o *Option[T]) Some() T {
	if !o.hasSome {
		panic("Someue is nil")
	}
	return o.val
}

func (o *Option[T]) IsSome() bool {
	return o.hasSome
}

func (o *Option[T]) UnwrapOr(def T) T {
	if o.hasSome {
		return o.val
	}
	return def
}

func (o *Option[T]) UnwrapOrF(def func() T) T {
	if o.hasSome {
		return o.val
	}
	return def()
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		val:     val,
		hasSome: true,
	}
}

func None[T any]() Option[T] {
	var o Option[T]
	return o
}

func Chain[T any](iters ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, i := range iters {
			for x := range i {
				if !yield(x) {
					return
				}
			}
		}
	}
}

func Map[T, K any](in iter.Seq[T], f func(T) K) iter.Seq[K] {
	return func(yield func(K) bool) {
		for i := range in {
			if !yield(f(i)) {
				return
			}
		}
	}
}

func Filter[T any](in iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i := range in {
			if pred(i) {
				if !yield(i) {
					return
				}
			}
		}
	}
}

func Reduce[T any](in iter.Seq[T], seed T, f func(a, b T) T) T {
	out := seed
	for v := range in {
		out = f(out, v)
	}
	return out
}

func Sum[T constraints.Integer | constraints.Float](in iter.Seq[T]) T {
	var zero T
	return Reduce(in, zero, func(a, b T) T {
		return a + b
	})
}
func StepRange[T constraints.Integer | constraints.Float](start, end, step T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for start < end {
			if !yield(start) {
				return
			}
			start += step
		}
	}
}

func Range[T constraints.Integer | constraints.Float](start, end T) iter.Seq[T] {
	return StepRange(start, end, 1)
}
