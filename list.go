package fn

import (
	"bytes"
	"errors"
	"fmt"
)

var IndexOutOfRange = errors.New("index out of range")

// A List is an immutable singly linked list that is safe for concurrent use
type List[T any] struct {
	next *List[T]
	val  T
}

// NewList creates and returns an new list with the given value at the first node
func NewList[T any](val T) *List[T] {
	return &List[T]{
		val: val,
	}
}

// Val returns the value stored at the current node in the list
func (l *List[T]) Val() any {
	return l.val
}

// Len returns the length of the list
func (l *List[T]) Len() int {
	i := 1
	y := l
	for !y.End() {
		i++
		y = y.next
	}

	return i
}

// String returns a string representation of the list
func (l *List[T]) String() string {
	if l == nil {
		return "nil"
	}
	b := bytes.NewBuffer(nil)
	b.WriteString("[")
	y := l
	for {
		b.WriteString(fmt.Sprintf("%v", y.val))
		if !y.End() {
			b.WriteString(", ")
		} else {
			break
		}
		y = y.next
	}
	b.WriteString("]")

	return b.String()
}

// End returns true if this is the end of the list
func (l *List[T]) End() bool {
	return l.next == nil
}

// Index returns the value stored at the given index if it exists
func (l *List[T]) Index(i int) (any, error) {
	x := 0
	y := l

	for x < i {

		if y.End() {
			return nil, IndexOutOfRange
		}
		y = y.next
		x++
	}

	return y.val, nil
}

// Prepend the given value onto a new list
func (l *List[T]) Prepend(val T) *List[T] {
	return &List[T]{
		next: l,
		val:  val,
	}
}

// Append the given value to the end of the list. This will reallocate the whole list
func (l *List[T]) Append(val T) *List[T] {
	// make a copy of this list
	n := &List[T]{}
	n.val = l.val

	//  if this is not the end, pass it down the line
	if !l.End() {
		n.next = l.next.Append(val)
	} else {
		n.next = &List[T]{
			val: val,
		}
	}

	return n
}

// Next returns the next node in the list
func (l *List[T]) Next() *List[T] {
	return l.next
}

func (l *List[T]) Each(f func(i T)) {
	if l == nil {
		return
	}

	f(l.val)
	l.Next().Each(f)
}

func (l *List[T]) Filter(f func(*List[T]) bool) *List[T] {
	if l == nil {
		return nil
	}

	if f(l) {

		n := NewList(l.val)
		n.next = n.next.Filter(f)
		return n
	}

	if l.End() {
		return nil
	}

	return l.next.Filter(f)
}
