package fn

import "iter"

func All(i iter.Seq[bool]) bool {
	one := false
	for x := range i {
		one = true
		if !x {
			return false
		}
	}
	return one
}

func Len[T any](i iter.Seq[T]) int {
	x := 0
	for range i {
		x++
	}
	return x
}
