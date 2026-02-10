package fn

import (
	"io/fs"
	"testing"

	"github.com/eliothedeman/check"
)

func TestRange(t *testing.T) {
	r := Range(0, 100)
	i := 0
	for v := range r {
		if i != v {
			t.Errorf("have %d want %d", v, i)
		}
		i++
	}
	if i != 100 {
		t.Error("Range should be 100 got ", i)
	}
}

func TestChain(t *testing.T) {
	i := Chain(Range(0, 2), Range(5, 10))

	sum := Sum(i)
	if sum != 36 {
		t.Error(sum)
	}
}

func TestResult(t *testing.T) {
	check.Eq(Sum(Map(Ok(100).Iter(), func(i int) float32 {
		return float32(i) + 20.1
	})), 120.1)
	check.Eq(Sum(Map(Err[int](fs.ErrNotExist).Iter(), func(i int) float32 {
		return float32(i) + 20.1
	})), 0)
}
