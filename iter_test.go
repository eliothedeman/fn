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

// --- Apply ---

func TestApply(t *testing.T) {
	doubled := slices.Collect(Apply(Range(1, 4), func(i int) int {
		return i * 2
	}))
	check.Eq(len(doubled), 3)
	check.Eq(doubled[0], 2)
	check.Eq(doubled[1], 4)
	check.Eq(doubled[2], 6)
}

func TestApplyTypeConversion(t *testing.T) {
	strs := slices.Collect(Apply(Range(0, 3), func(i int) string {
		return string(rune('a' + i))
	}))
	check.Eq(len(strs), 3)
	check.Eq(strs[0], "a")
	check.Eq(strs[1], "b")
	check.Eq(strs[2], "c")
}

func TestApplyEmpty(t *testing.T) {
	vals := slices.Collect(Apply(Range(0, 0), func(i int) int {
		return i * 2
	}))
	check.Eq(len(vals), 0)
}

func TestApplyEarlyBreak(t *testing.T) {
	count := 0
	for range Apply(Range(0, 100), func(i int) int {
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
	check.Eq(HasValue(r), true)
	val, err := Unpack(r)
	check.Eq(val, 42)
	check.Nil(err)
}

func TestErr(t *testing.T) {
	r := Err[int](fs.ErrNotExist)
	val, err := Unpack(r)
	check.Eq(val, 0)
	check.NotNil(err)
	check.ErrIs(err, fs.ErrNotExist)
}

func TestTry(t *testing.T) {
	r := Try(10, nil)
	check.Eq(HasValue(r), true)
	val, err := Unpack(r)
	check.Eq(val, 10)
	check.Nil(err)
}

func TestTryWithError(t *testing.T) {
	r := Try(0, fs.ErrPermission)
	val, err := Unpack(r)
	check.Eq(val, 0)
	check.NotNil(err)
	check.ErrIs(err, fs.ErrPermission)
}

func TestResultUnpack(t *testing.T) {
	val, err := Unpack(Ok("hello"))
	check.Eq(val, "hello")
	check.Nil(err)

	val2, err2 := Unpack(Err[string](fs.ErrClosed))
	check.Eq(val2, "")
	check.NotNil(err2)
}

func TestResultIterOk(t *testing.T) {
	vals := slices.Collect(Iter(Ok(99)))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 99)
}

func TestResultIterErr(t *testing.T) {
	vals := slices.Collect(Iter(Err[int](fs.ErrNotExist)))
	check.Eq(len(vals), 0)
}

func TestResultIterErrMethod(t *testing.T) {
	r := Err[int](fs.ErrNotExist)
	errs := slices.Collect(IterErr(r))
	check.Eq(len(errs), 1)
	check.ErrIs(errs[0], fs.ErrNotExist)
}

func TestResultIterErrOnOk(t *testing.T) {
	r := Ok(10)
	errs := slices.Collect(IterErr(r))
	check.Eq(len(errs), 0)
}

func TestResultUnwrap(t *testing.T) {
	check.Eq(Unwrap(Ok(42)), 42)
}

func TestResultUnwrapPanicsOnErr(t *testing.T) {
	check.Panics(func() {
		Unwrap(Err[int](errors.New("fail")))
	})
}

func TestResultUnwrapOr(t *testing.T) {
	check.Eq(UnwrapOr(Ok(10), 99), 10)
	check.Eq(UnwrapOr(Err[int](errors.New("fail")), 99), 99)
}

func TestResultUnwrapOrF(t *testing.T) {
	called := false
	val := UnwrapOrF(Ok(10), func() int {
		called = true
		return 99
	})
	check.Eq(val, 10)
	check.Eq(called, false)

	val2 := UnwrapOrF(Err[int](errors.New("fail")), func() int {
		return 99
	})
	check.Eq(val2, 99)
}

func TestResultHasValue(t *testing.T) {
	check.Eq(HasValue(Ok(1)), true)
	check.Eq(HasValue(Err[int](errors.New("fail"))), false)
}

func TestResultIsEmpty(t *testing.T) {
	check.Eq(IsEmpty(Err[int](errors.New("fail"))), true)
	check.Eq(IsEmpty(Ok(1)), false)
}

func TestResultApplyViaIter(t *testing.T) {
	vals := slices.Collect(Apply(Iter(Ok(5)), func(i int) int {
		return i * 10
	}))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 50)

	vals2 := slices.Collect(Apply(Iter(Err[int](fs.ErrNotExist)), func(i int) int {
		return i * 10
	}))
	check.Eq(len(vals2), 0)
}

func TestResultSumViaIter(t *testing.T) {
	check.Eq(Sum(Iter(Ok(42))), 42)
	check.Eq(Sum(Iter(Err[int](fs.ErrNotExist))), 0)
}

func TestResultChainIters(t *testing.T) {
	sum := Sum(Chain(Iter(Ok(10)), Iter(Ok(20)), Iter(Err[int](fs.ErrNotExist)), Iter(Ok(30))))
	check.Eq(sum, 60)
}

// --- Option ---

func TestSome(t *testing.T) {
	o := Some(42)
	check.Eq(HasValue(o), true)
	check.Eq(Unwrap(o), 42)
}

func TestNone(t *testing.T) {
	o := None[int]()
	check.Eq(HasValue(o), false)
	check.Eq(IsEmpty(o), true)
}

func TestOptionUnwrapPanicsOnNone(t *testing.T) {
	check.Panics(func() {
		Unwrap(None[int]())
	})
}

func TestOptionUnwrapOr(t *testing.T) {
	check.Eq(UnwrapOr(Some(10), 99), 10)
	check.Eq(UnwrapOr(None[int](), 99), 99)
}

func TestOptionUnwrapOrF(t *testing.T) {
	called := false
	val := UnwrapOrF(Some(10), func() int {
		called = true
		return 99
	})
	check.Eq(val, 10)
	check.Eq(called, false)

	val2 := UnwrapOrF(None[int](), func() int {
		return 99
	})
	check.Eq(val2, 99)
}

func TestOptionIter(t *testing.T) {
	vals := slices.Collect(Iter(Some(7)))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 7)
}

func TestOptionIterNone(t *testing.T) {
	vals := slices.Collect(Iter(None[int]()))
	check.Eq(len(vals), 0)
}

func TestOptionIterString(t *testing.T) {
	vals := slices.Collect(Iter(Some("hello")))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], "hello")
}

func TestOptionApplyViaIter(t *testing.T) {
	vals := slices.Collect(Apply(Iter(Some(5)), func(i int) int {
		return i * 3
	}))
	check.Eq(len(vals), 1)
	check.Eq(vals[0], 15)
}

func TestOptionSumViaIter(t *testing.T) {
	check.Eq(Sum(Iter(Some(42))), 42)
	check.Eq(Sum(Iter(None[int]())), 0)
}

func TestOptionChainIters(t *testing.T) {
	sum := Sum(Chain(Iter(Some(1)), Iter(None[int]()), Iter(Some(2))))
	check.Eq(sum, 3)
}

func TestOptionUnwrapOrZeroValue(t *testing.T) {
	check.Eq(UnwrapOr(None[string](), ""), "")
	check.Eq(UnwrapOr(None[int](), 0), 0)
}

// --- Composition ---

func TestFilterApplyReduceComposition(t *testing.T) {
	// Sum of doubled even numbers from 0..10
	result := Reduce(
		Apply(
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

func TestApplyChain(t *testing.T) {
	a := Apply(Range(0, 3), func(i int) int { return i * 10 })
	b := Apply(Range(0, 3), func(i int) int { return i * 100 })
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
	sum := Sum(Chain(Iter(Ok(10)), Iter(Some(20))))
	check.Eq(sum, 30)
}

func TestZip(t *testing.T) {
	for a, b := range Zip(Range(0, 10), Range(0, 10)) {
		check.Eq(a, b)
	}
	for a, b := range Zip(Range(0, 10), Range(0, 11)) {
		check.Eq(a, b)
	}
	for a, b := range Zip(Range(0, 9), Range(0, 11)) {
		check.Eq(a, b)
	}
}
