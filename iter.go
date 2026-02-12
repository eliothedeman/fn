// Package fn provides functional programming primitives built on Go's iter.Seq
// iterators. It uses a lisp-style free function API: container types like
// [Result] and [Option] expose only an Iter method, while all inspection and
// extraction is done through shared free functions like [Unwrap], [HasValue],
// and [UnwrapOr]. Iterator combinators ([Apply], [Filter], [Reduce], [Chain])
// compose freely with these container iterators to build data pipelines.
package fn

import (
	"iter"
	"slices"

	"golang.org/x/exp/constraints"
)

// Iterable is the shared interface that connects container types to the
// iterator combinator pipeline. Both [Result] and [Option] satisfy Iterable,
// so they can be passed directly to [Iter] and from there into [Apply], [Filter],
// [Chain], [Reduce], and any other function that accepts an iter.Seq.
type Iterable[T any] interface {
	Iter() iter.Seq[T]
}

// Iter extracts an iter.Seq from any [Iterable] container. This is the
// bridge between value containers and the iterator pipeline—use it to feed
// a [Result] or [Option] into combinators like [Apply] or [Chain]:
//
//	sum := fn.Sum(fn.Chain(fn.Iter(fn.Ok(10)), fn.Iter(fn.Some(20))))
func Iter[T any](x Iterable[T]) iter.Seq[T] {
	return x.Iter()
}

type unwrappable[T any] interface {
	unwrap() (T, bool)
}

// Unwrap extracts the value from a [Result] or [Option], panicking if the
// container is empty (Err or None). Use this when the caller can guarantee the
// container holds a value and a missing value represents a programming error.
// Prefer [UnwrapOr] or [UnwrapOrF] when a fallback is more appropriate.
func Unwrap[T any](x unwrappable[T]) T {
	val, ok := x.unwrap()
	if !ok {
		panic("called Unwrap on an empty value")
	}
	return val
}

// UnwrapOr extracts the value from a [Result] or [Option], returning def if
// the container is empty. This is the simplest safe extraction when you have a
// reasonable default:
//
//	name := fn.UnwrapOr(fn.None[string](), "anonymous")
func UnwrapOr[T any](x unwrappable[T], def T) T {
	val, ok := x.unwrap()
	if ok {
		return val
	}
	return def
}

// UnwrapOrF extracts the value from a [Result] or [Option], calling def to
// produce a fallback only when the container is empty. Use this instead of
// [UnwrapOr] when computing the default is expensive or has side effects,
// since the function is not called on the happy path.
func UnwrapOrF[T any](x unwrappable[T], def func() T) T {
	val, ok := x.unwrap()
	if ok {
		return val
	}
	return def()
}

// HasValue reports whether a [Result] or [Option] contains a value (Ok or
// Some). It is the complement of [IsEmpty] and works uniformly across both
// container types, so callers don't need to remember IsOk vs IsSome.
func HasValue[T any](x unwrappable[T]) bool {
	_, ok := x.unwrap()
	return ok
}

// IsEmpty reports whether a [Result] or [Option] is empty (Err or None).
// It is the complement of [HasValue].
func IsEmpty[T any](x unwrappable[T]) bool {
	_, ok := x.unwrap()
	return !ok
}

// Result holds either a value of type T or an error—never both. Use [Ok] and
// [Err] to construct results, [Try] to capture a standard Go (T, error) return,
// and the shared free functions ([Unwrap], [UnwrapOr], [HasValue], [Iter]) to
// inspect or extract the value. Result satisfies [Iterable], yielding the value
// on Ok and nothing on Err, so it plugs directly into iterator pipelines:
//
//	total := fn.Sum(fn.Chain(fn.Iter(fn.Try(strconv.Atoi("3"))), fn.Range(0, 10)))
type Result[T any] struct {
	val T
	err error
}

// Try constructs a [Result] from a (T, error) pair, the standard Go multi-return
// convention. Wrap any stdlib or third-party call to lift it into the fn pipeline:
//
//	r := fn.Try(os.Open("config.json"))
func Try[T any](t T, err error) Result[T] {
	return Result[T]{val: t, err: err}
}

// Ok constructs a successful [Result] containing val. The result carries no
// error, so [HasValue] returns true and [Iter] yields val.
func Ok[T any](t T) Result[T] {
	return Result[T]{val: t}
}

// Err constructs a failed [Result] carrying err and the zero value of T.
// [HasValue] returns false, and [Iter] yields nothing, effectively filtering
// this result out of any iterator pipeline it participates in.
func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}

// Unpack destructures a [Result] back into Go's conventional (T, error) pair.
// Use this at API boundaries where you need to return or switch on the error:
//
//	val, err := fn.Unpack(fn.Try(strconv.Atoi(input)))
//	if err != nil { ... }
func Unpack[T any](r Result[T]) (T, error) {
	return r.val, r.err
}

func (r Result[T]) unwrap() (T, bool) {
	return r.val, r.err == nil
}

// Iter returns an iterator that yields the contained value if the Result is Ok,
// or nothing if it is Err. This is the method that satisfies [Iterable] and
// allows Results to participate in iterator pipelines.
func (r Result[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if r.err == nil {
			yield(r.val)
		}
	}
}

// IterErr returns an iterator over the error side of a [Result]. It yields
// the error if the Result is Err, or nothing if it is Ok. This is useful for
// collecting or inspecting errors from a set of results:
//
//	for e := range fn.IterErr(result) {
//	    log.Println("operation failed:", e)
//	}
func IterErr[T any](r Result[T]) iter.Seq[error] {
	return func(yield func(error) bool) {
		if r.err != nil {
			yield(r.err)
		}
	}
}

// Option represents a value that may or may not be present, without using nil
// pointers or sentinel values. Use [Some] and [None] to construct options, and
// the shared free functions ([Unwrap], [UnwrapOr], [HasValue], [Iter]) to
// inspect or extract the value. Option satisfies [Iterable], yielding the value
// on Some and nothing on None, so it composes with iterator pipelines just like
// [Result]:
//
//	sum := fn.Sum(fn.Chain(fn.Iter(fn.Some(1)), fn.Iter(fn.None[int]()), fn.Iter(fn.Some(2))))
//	// sum == 3
type Option[T any] struct {
	val     T
	hasSome bool
}

func (o Option[T]) unwrap() (T, bool) {
	return o.val, o.hasSome
}

// Iter returns an iterator that yields the contained value if the Option is
// Some, or nothing if it is None. This is the method that satisfies [Iterable]
// and allows Options to participate in iterator pipelines.
func (o Option[T]) Iter() iter.Seq[T] {
	return func(yield func(T) bool) {
		if o.hasSome {
			yield(o.val)
		}
	}
}

// Some constructs an [Option] that contains val. [HasValue] returns true and
// [Iter] yields val.
func Some[T any](val T) Option[T] {
	return Option[T]{
		val:     val,
		hasSome: true,
	}
}

// None constructs an empty [Option] with no value. [HasValue] returns false
// and [Iter] yields nothing, so a None is silently skipped in iterator
// pipelines.
func None[T any]() Option[T] {
	var o Option[T]
	return o
}

// Chain concatenates multiple iterators into a single iterator that yields
// all elements from the first, then all from the second, and so on. Use it to
// merge sequences from different sources—including [Result] and [Option]
// iterators—into one pipeline:
//
//	all := fn.Chain(fn.Range(0, 5), fn.Iter(fn.Ok(99)), fn.Iter(fn.Some(100)))
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

// Zip two iterators together
func Zip[T any, K any](a iter.Seq[T], b iter.Seq[K]) iter.Seq2[T, K] {
	return func(yield func(T, K) bool) {
		pa, stopa := iter.Pull(a)
		pb, stopb := iter.Pull(b)
		for {
			ax, aok := pa()
			bx, bok := pb()
			if (!aok || !bok) || !yield(ax, bx) {
				stopa()
				stopb()
				return
			}
		}
	}
}

// Apply transforms each element in an iterator by applying f, producing a new
// iterator of the mapped type. The transformation is lazy—f is called only
// as elements are consumed:
//
//	doubled := fn.Apply(fn.Range(1, 4), func(i int) int { return i * 2 })
//	// yields 2, 4, 6
func Apply[T, K any](in iter.Seq[T], f func(T) K) iter.Seq[K] {
	return func(yield func(K) bool) {
		for i := range in {
			if !yield(f(i)) {
				return
			}
		}
	}
}

// Filter produces an iterator that yields only the elements for which pred
// returns true. Like [Apply], evaluation is lazy—pred is called as elements are
// consumed, and the pipeline short-circuits on early break:
//
//	evens := fn.Filter(fn.Range(0, 10), func(i int) bool { return i%2 == 0 })
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

// Reduce folds an iterator into a single value by repeatedly applying f to an
// accumulator (starting at seed) and each successive element. This is the
// primary way to collapse an iterator pipeline into a scalar result:
//
//	product := fn.Reduce(fn.Range(1, 6), 1, func(a, b int) int { return a * b })
//	// product == 120
func Reduce[T any](in iter.Seq[T], seed T, f func(a, b T) T) T {
	out := seed
	for v := range in {
		out = f(out, v)
	}
	return out
}

// Sum is a convenience specialization of [Reduce] that adds all numeric values
// in an iterator. It returns the zero value for an empty sequence:
//
//	fn.Sum(fn.Range(0, 100)) // 4950
func Sum[T constraints.Integer | constraints.Float](in iter.Seq[T]) T {
	var zero T
	return Reduce(in, zero, func(a, b T) T {
		return a + b
	})
}

// StepRange produces an iterator of numeric values from start (inclusive) to
// end (exclusive), advancing by step each time. If start >= end the iterator
// is empty. Use [Range] for the common case where step is 1.
//
//	fn.StepRange(0, 10, 3) // yields 0, 3, 6, 9
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

// Range produces an iterator of numeric values from start (inclusive) to end
// (exclusive) with a step of 1. It works with any integer or float type.
// Range is often the starting point for a pipeline:
//
//	fn.Sum(fn.Filter(fn.Range(0, 20), func(i int) bool { return i%2 == 0 }))
func Range[T constraints.Integer | constraints.Float](start, end T) iter.Seq[T] {
	return StepRange(start, end, 1)
}

// Reverse an iterator such that it iterates back to front
func Reverse[T any](i iter.Seq[T]) iter.Seq[T] {
	s := slices.Collect(i)
	slices.Reverse(s)
	return slices.Values(s)
}

func Enumerate[T any](x iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for y := range x {
			if !yield(i, y) {
				return
			}
			i++
		}
	}
}
