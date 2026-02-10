package fn

import (
	"errors"
	"io/fs"
	"slices"
	"testing"

	"github.com/eliothedeman/check"
)

// --- Range / StepRange ---

func TestRange(t *testing.T) {
	vals := slices.Collect(Range(0, 5))
	check.Eq(len(vals), 5)
	check.Eq(vals[0], 0)
	check.Eq(vals[4], 4)
}

func TestRangeEmpty(t *testing.T) {
	vals := slices.Collect(Range(5, 5))
	check.Eq(len(vals), 0)
}

func TestRangeNoNegativeProgress(t *testing.T) {
	vals := slices.Collect(Range(10, 5))
	check.Eq(len(vals), 0)
}

func TestRangeLargeSequence(t *testing.T) {
	check.Eq(Sum(Range(0, 1000)), 499500)
}

func TestRangeFloat(t *testing.T) {
	vals := slices.Collect(Range(0.0, 3.0))
	check.Eq(len(vals), 3)
	check.Eq(vals[0], 0.0)
	check.Eq(vals[1], 1.0)
	check.Eq(vals[2], 2.0)
}

func TestStepRange(t *testing.T) {
	vals := slices.Collect(StepRange(0, 10, 2))
	check.Eq(len(vals), 5)
	check.Eq(vals[0], 0)
	check.Eq(vals[1], 2)
	check.Eq(vals[4], 8)
}

func TestStepRangeFloat(t *testing.T) {
	vals := slices.Collect(StepRange(0.0, 1.0, 0.25))
	check.Eq(len(vals), 4)
	check.Eq(vals[0], 0.0)
	check.Eq(vals[3], 0.75)
}

func TestStepRangeEmpty(t *testing.T) {
	vals := slices.Collect(StepRange(10, 5, 1))
	check.Eq(len(vals), 0)
}

func TestStepRangeLargeStep(t *testing.T) {
	vals := slices.Collect(StepRange(0, 10, 100))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 0)
}

func TestRangeEarlyBreak(t *testing.T) {
	count := 0
	for v := range Range(0, 100) {
		count++
		if v == 2 {
			break
		}
	}
	check.Eq(count, 3)
}

// --- Chain ---

func TestChain(t *testing.T) {
	vals := slices.Collect(Chain(Range(0, 3), Range(10, 13)))
	check.Eq(len(vals), 6)
	check.Eq(vals[0], 0)
	check.Eq(vals[2], 2)
	check.Eq(vals[3], 10)
	check.Eq(vals[5], 12)
}

func TestChainEmpty(t *testing.T) {
	vals := slices.Collect(Chain[int]())
	check.Eq(len(vals), 0)
}

func TestChainSingle(t *testing.T) {
	vals := slices.Collect(Chain(Range(0, 3)))
	check.Eq(len(vals), 3)
	check.Eq(vals[0], 0)
	check.Eq(vals[2], 2)
}

func TestChainMultiple(t *testing.T) {
	sum := Sum(Chain(Range(0, 5), Range(5, 10), Range(10, 15)))
	check.Eq(sum, 105)
}

func TestChainEarlyBreak(t *testing.T) {
	count := 0
	for range Chain(Range(0, 100), Range(0, 100)) {
		count++
		if count == 5 {
			break
		}
	}
	check.Eq(count, 5)
}

// --- Map ---

func TestMap(t *testing.T) {
	doubled := slices.Collect(Map(Range(1, 4), func(i int) int {
		return i * 2
	}))
	check.Eq(len(doubled), 3)
	check.Eq(doubled[0], 2)
	check.Eq(doubled[1], 4)
	check.Eq(doubled[2], 6)
}

func TestMapTypeConversion(t *testing.T) {
	strs := slices.Collect(Map(Range(0, 3), func(i int) string {
		return string(rune('a' + i))
	}))
	check.Eq(len(strs), 3)
	check.Eq(strs[0], "a")
	check.Eq(strs[1], "b")
	check.Eq(strs[2], "c")
}

func TestMapEmpty(t *testing.T) {
	vals := slices.Collect(Map(Range(0, 0), func(i int) int {
		return i * 2
	}))
	check.Eq(len(vals), 0)
}

func TestMapEarlyBreak(t *testing.T) {
	count := 0
	for range Map(Range(0, 100), func(i int) int {
		count++
		return i
	}) {
		if count == 3 {
			break
		}
	}
	check.Eq(count, 3)
}

// --- Filter ---

func TestFilter(t *testing.T) {
	evens := slices.Collect(Filter(Range(0, 10), func(i int) bool {
		return i%2 == 0
	}))
	check.Eq(len(evens), 5)
	check.Eq(evens[0], 0)
	check.Eq(evens[1], 2)
	check.Eq(evens[4], 8)
}

func TestFilterNoneMatch(t *testing.T) {
	vals := slices.Collect(Filter(Range(0, 5), func(i int) bool {
		return i > 100
	}))
	check.Eq(len(vals), 0)
}

func TestFilterAllMatch(t *testing.T) {
	vals := slices.Collect(Filter(Range(0, 5), func(i int) bool {
		return true
	}))
	check.Eq(len(vals), 5)
}

func TestFilterEarlyBreak(t *testing.T) {
	count := 0
	for range Filter(Range(0, 100), func(i int) bool {
		return i%2 == 0
	}) {
		count++
		if count == 3 {
			break
		}
	}
	check.Eq(count, 3)
}

// --- Reduce ---

func TestReduce(t *testing.T) {
	sum := Reduce(Range(1, 5), 0, func(a, b int) int {
		return a + b
	})
	check.Eq(sum, 10)
}

func TestReduceProduct(t *testing.T) {
	product := Reduce(Range(1, 6), 1, func(a, b int) int {
		return a * b
	})
	check.Eq(product, 120)
}

func TestReduceEmpty(t *testing.T) {
	result := Reduce(Range(0, 0), 42, func(a, b int) int {
		return a + b
	})
	check.Eq(result, 42)
}

func TestReduceMax(t *testing.T) {
	maxVal := Reduce(Chain(Range(0, 5), Range(3, 8)), 0, func(a, b int) int {
		if b > a {
			return b
		}
		return a
	})
	check.Eq(maxVal, 7)
}

// --- Sum ---

func TestSum(t *testing.T) {
	check.Eq(Sum(Range(0, 100)), 4950)
}

func TestSumEmpty(t *testing.T) {
	check.Eq(Sum(Range(0, 0)), 0)
}

func TestSumSingleElement(t *testing.T) {
	check.Eq(Sum(Range(5, 6)), 5)
}

func TestSumFloat(t *testing.T) {
	check.Eq(Sum(Range(0.0, 4.0)), 6.0)
}

// --- Result ---

func TestOk(t *testing.T) {
	r := Ok(42)
	check.Eq(r.IsOk(), true)
	val, err := r.Unpack()
	check.Eq(val, 42)
	check.Nil(err)
}

func TestErr(t *testing.T) {
	r := Err[int](fs.ErrNotExist)
	val, err := r.Unpack()
	check.Eq(val, 0)
	check.NotNil(err)
	check.ErrIs(err, fs.ErrNotExist)
}

func TestTry(t *testing.T) {
	r := Try(10, nil)
	check.Eq(r.IsOk(), true)
	val, err := r.Unpack()
	check.Eq(val, 10)
	check.Nil(err)
}

func TestTryWithError(t *testing.T) {
	r := Try(0, fs.ErrPermission)
	val, err := r.Unpack()
	check.Eq(val, 0)
	check.NotNil(err)
	check.ErrIs(err, fs.ErrPermission)
}

func TestResultUnpack(t *testing.T) {
	val, err := Ok("hello").Unpack()
	check.Eq(val, "hello")
	check.Nil(err)

	val2, err2 := Err[string](fs.ErrClosed).Unpack()
	check.Eq(val2, "")
	check.NotNil(err2)
}

func TestResultIterOk(t *testing.T) {
	vals := slices.Collect(Ok(99).Iter())
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 99)
}

func TestResultIterErr(t *testing.T) {
	vals := slices.Collect(Err[int](fs.ErrNotExist).Iter())
	check.Eq(len(vals), 0)
}

func TestResultIterErrMethod(t *testing.T) {
	r := Err[int](fs.ErrNotExist)
	errs := slices.Collect(r.IterErr())
	check.Eq(len(errs), 1)
	check.ErrIs(errs[0], fs.ErrNotExist)
}

func TestResultIterErrOnOk(t *testing.T) {
	r := Ok(10)
	errs := slices.Collect(r.IterErr())
	check.Eq(len(errs), 0)
}

func TestResultIsOk(t *testing.T) {
	check.Eq(Ok(1).IsOk(), true)
	check.Eq(Err[int](errors.New("fail")).IsOk(), false)
}

func TestResultMapViaIter(t *testing.T) {
	vals := slices.Collect(Map(Ok(5).Iter(), func(i int) int {
		return i * 10
	}))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 50)

	vals2 := slices.Collect(Map(Err[int](fs.ErrNotExist).Iter(), func(i int) int {
		return i * 10
	}))
	check.Eq(len(vals2), 0)
}

func TestResultSumViaIter(t *testing.T) {
	check.Eq(Sum(Ok(42).Iter()), 42)
	check.Eq(Sum(Err[int](fs.ErrNotExist).Iter()), 0)
}

func TestResultChainIters(t *testing.T) {
	sum := Sum(Chain(Ok(10).Iter(), Ok(20).Iter(), Err[int](fs.ErrNotExist).Iter(), Ok(30).Iter()))
	check.Eq(sum, 60)
}

// --- Option ---

func TestSome(t *testing.T) {
	o := Some(42)
	check.Eq(o.IsSome(), true)
	check.Eq(o.Some(), 42)
}

func TestNone(t *testing.T) {
	o := None[int]()
	check.Eq(o.IsSome(), false)
}

func TestOptionSomePanicsOnNone(t *testing.T) {
	check.Panics(func() {
		n := None[int]()
		n.Some()
	})
}

func TestOptionUnwrapOr(t *testing.T) {
	s := Some(10)
	check.Eq(s.UnwrapOr(99), 10)
	n := None[int]()
	check.Eq(n.UnwrapOr(99), 99)
}

func TestOptionUnwrapOrF(t *testing.T) {
	called := false
	s := Some(10)
	val := s.UnwrapOrF(func() int {
		called = true
		return 99
	})
	check.Eq(val, 10)
	check.Eq(called, false)

	n := None[int]()
	val2 := n.UnwrapOrF(func() int {
		return 99
	})
	check.Eq(val2, 99)
}

func TestOptionIter(t *testing.T) {
	s := Some(7)
	vals := slices.Collect(s.Iter())
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 7)
}

func TestOptionIterNone(t *testing.T) {
	n := None[int]()
	vals := slices.Collect(n.Iter())
	check.Eq(len(vals), 0)
}

func TestOptionIterString(t *testing.T) {
	s := Some("hello")
	vals := slices.Collect(s.Iter())
	check.Eq(len(vals), 1)
	check.Eq(vals[0], "hello")
}

func TestOptionMapViaIter(t *testing.T) {
	s := Some(5)
	vals := slices.Collect(Map(s.Iter(), func(i int) int {
		return i * 3
	}))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 15)
}

func TestOptionSumViaIter(t *testing.T) {
	s := Some(42)
	check.Eq(Sum(s.Iter()), 42)
	n := None[int]()
	check.Eq(Sum(n.Iter()), 0)
}

func TestOptionChainIters(t *testing.T) {
	s1 := Some(1)
	n := None[int]()
	s2 := Some(2)
	sum := Sum(Chain(s1.Iter(), n.Iter(), s2.Iter()))
	check.Eq(sum, 3)
}

func TestOptionUnwrapOrZeroValue(t *testing.T) {
	ns := None[string]()
	check.Eq(ns.UnwrapOr(""), "")
	ni := None[int]()
	check.Eq(ni.UnwrapOr(0), 0)
}

// --- Composition ---

func TestFilterMapReduceComposition(t *testing.T) {
	// Sum of doubled even numbers from 0..10
	result := Reduce(
		Map(
			Filter(Range(0, 10), func(i int) bool { return i%2 == 0 }),
			func(i int) int { return i * 2 },
		),
		0,
		func(a, b int) int { return a + b },
	)
	check.Eq(result, 40) // (0+2+4+6+8)*2 = 40
}

func TestChainFilterSum(t *testing.T) {
	sum := Sum(Filter(
		Chain(Range(0, 5), Range(10, 15)),
		func(i int) bool { return i%2 == 0 },
	))
	// evens: 0, 2, 4, 10, 12, 14
	check.Eq(sum, 42)
}

func TestMapChain(t *testing.T) {
	a := Map(Range(0, 3), func(i int) int { return i * 10 })
	b := Map(Range(0, 3), func(i int) int { return i * 100 })
	vals := slices.Collect(Chain(a, b))
	check.Eq(len(vals), 6)
	check.Eq(vals[0], 0)
	check.Eq(vals[1], 10)
	check.Eq(vals[2], 20)
	check.Eq(vals[3], 0)
	check.Eq(vals[4], 100)
	check.Eq(vals[5], 200)
}

func TestResultOptionInterop(t *testing.T) {
	// Chain a Result iter with an Option iter
	s := Some(20)
	sum := Sum(Chain(Ok(10).Iter(), s.Iter()))
	check.Eq(sum, 30)
}
