package vec

import (
	"math"
	"testing"
)

func Test_Ops(t *testing.T) {
	a := V3{1, 2, 3}
	b := V3{4, 5, 6}

	good_sum := V3{5, 7, 9}
	if a.Sum(b) != good_sum {
		t.Error("FAIL")
	}

	good_sub := V3{3, 3, 3}
	if b.Sub(a) != good_sub {
		t.Error("FAIL")
	}

	good_cross := V3{-3, 6, -3}
	if a.Cross(b) != good_cross {
		t.Error("FAIL")
	}

	good_scale := V3{2, 4, 6}
	if a.Scale(2) != good_scale {
		t.Error("FAIL")
	}

	good_len := float32(math.Sqrt(1 + 4 + 9))
	if a.Length() != good_len {
		t.Error("FAIL")
	}

	good_dot := float32(4 + 10 + 18)
	if a.Dot(b) != good_dot {
		t.Error("FAIL")
	}

	l := float32(math.Sqrt(1 + 4 + 9))
	good_normalize := V3{a[0] / l, a[1] / l, a[2] / l}
	if a.Normalize() != good_normalize {
		t.Error("FAIL")
	}

}
