package fn

import "constraints"

type Option[T any] struct {
	val    T
	hasVal bool
}

func (o *Option[T]) Val() T {
	if !o.hasVal {
		panic("Value is nil")
	}
	return o.val
}

func (o *Option[T]) HasVal() bool {
	return o.hasVal
}

func (o *Option[T]) ValOr(def T) T {
	if o.hasVal {
		return o.val
	}
	return def
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		val:    val,
		hasVal: true,
	}
}

func None[T any]() Option[T] {
	var o Option[T]
	return o
}

type Iter[T any] struct {
	next    func() Option[T]
	nextVal T
}

func (i *Iter[T]) Next() bool {
	nv := i.next()
	if nv.hasVal {
		i.nextVal = nv.val
		return true
	}
	return false
}

func (i *Iter[T]) Val() T {
	return i.nextVal
}

func (i *Iter[T]) Collect() []T {
	var out []T
	for i.Next() {
		out = append(out, i.Val())
	}
	return out
}

func Range[T constraints.Integer](start, end T) *Iter[T] {
	var i T = start
	return &Iter[T]{
		next: func() Option[T] {
			if i < end {
				out := Some(i)
				i++
				return out
			}
			return None[T]()
		},
		nextVal: start,
	}
}

func Map[T, K any](in *Iter[T], f func(T) K) *Iter[K] {
	i := new(Iter[K])
	i.next = func() Option[K] {
		o := in.next()
		if o.HasVal() {
			return Some(f(o.val))
		}
		return None[K]()
	}
	return i
}

func Filter[T any](in *Iter[T], pred func(T) bool) *Iter[T] {
	i := new(Iter[T])
	i.next = func() Option[T] {
		for {
			o := in.next()
			if o.HasVal() {
				if pred(o.val) {
					return o
				}
			} else {
				return None[T]()
			}
		}
	}
	return i
}
