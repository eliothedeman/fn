package fn

import (
	"testing"

	"github.com/eliothedeman/check"
)

func TestVecReverse(t *testing.T) {
	x := Vec[int]{1, 2, 3}
	check.SliceEq(Collect(Reverse(x.Iter())), Vec[int]{3, 2, 1}, "Incorrect reverse")
}
