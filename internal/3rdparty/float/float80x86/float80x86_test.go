package float80x86

import (
	"math"
	"math/big"
	"testing"
)

func TestNewFromBits(t *testing.T) {
	golden := []struct {
		se   uint16
		m    uint64
		want float64
	}{
		// Special numbers.
		// 0 111111111111111 10 non-zero = +NaN
		{se: 0x7FFF, m: 0xBFFFFFFFFFFFFFFF, want: math.NaN()},
		// -NaN
		// 1 111111111111111 10 non-zero = -NaN
		{se: 0xFFFF, m: 0xBFFFFFFFFFFFFFFF, want: -math.NaN()},

		// from: https://docs.oracle.com/cd/E19957-01/806-3568/ncg_math.html#960

		// 0000 00000000 00000000 = 0.0
		{se: 0x0000, m: 0x0000000000000000, want: 0.0},
		// 8000 00000000 00000000 = -0.0
		{se: 0x8000, m: 0x0000000000000000, want: math.Copysign(0, -1)},
		// 3FFF 80000000 00000000 = 1.0
		{se: 0x3FFF, m: 0x8000000000000000, want: 1.0},
		// 4000 80000000 00000000 = 2.0
		{se: 0x4000, m: 0x8000000000000000, want: 2.0},
		// 7FFE FFFFFFFF FFFFFFFF = 1.18973149535723176505e+4932 (max normal)
		//{se: 0x7FFE, m: 0xFFFFFFFFFFFFFFFF, want: 1.18973149535723176505e+4932},
		// 0001 80000000 00000000 = 3.36210314311209350626e-4932 (min positive normal)
		//{se: 0x0001, m: 0x8000000000000000, want: 3.36210314311209350626e-4932},
		// 0000 7FFFFFFF FFFFFFFF = 3.36210314311209350608e-4932 (max subnormal)
		//{se: 0x0000, m: 0x7FFFFFFFFFFFFFFF, want: 3.36210314311209350608e-4932},
		// 0000 00000000 00000001 = 3.64519953188247460253e-4951 (min positive subnormal)
		//{se: 0x0000, m: 0x0000000000000001, want: 3.64519953188247460253e-4951},
		// 7FFF 80000000 00000000 = infinity
		{se: 0x7FFF, m: 0x8000000000000000, want: math.Inf(1)},
		// FFFF 80000000 00000000 = -infinity
		{se: 0xFFFF, m: 0x8000000000000000, want: math.Inf(-1)},

		// 2^i
		// TODO: add test cases for 2^i
	}
	for _, g := range golden {
		f := NewFromBits(g.se, g.m)
		got, _ := f.Float64()
		wantBits := math.Float64bits(g.want)
		gotBits := math.Float64bits(got)
		//fmt.Printf("bits: 0x%04X (%v)\n", g.bits, g.want)
		if wantBits != gotBits {
			t.Errorf("0x%04X %016X: number mismatch; expected 0x%016X (%v), got 0x%016X (%v)", g.se, g.m, wantBits, g.want, gotBits, got)
		}
	}
}

func TestNewFromFloat64(t *testing.T) {
	golden := []struct {
		in  float64
		se  uint16
		m   uint64
		acc big.Accuracy
	}{
		// Special numbers.
		// 0 111111111111111 10 non-zero = +NaN
		{in: math.NaN(), se: 0x7FFF, m: 0xBFFFFFFFFFFFFFFF, acc: big.Exact},
		// -NaN
		// 1 111111111111111 10 non-zero = -NaN
		{in: -math.NaN(), se: 0xFFFF, m: 0xBFFFFFFFFFFFFFFF, acc: big.Exact},

		// from: https://docs.oracle.com/cd/E19957-01/806-3568/ncg_math.html#960

		// 0000 00000000 00000000 = 0.0
		{in: 0.0, se: 0x0000, m: 0x0000000000000000, acc: big.Exact},
		// 8000 00000000 00000000 = -0.0
		{in: math.Copysign(0, -1), se: 0x8000, m: 0x0000000000000000, acc: big.Exact},
		// 3FFF 80000000 00000000 = 1.0
		{in: 1.0, se: 0x3FFF, m: 0x8000000000000000, acc: big.Exact},
		// 4000 80000000 00000000 = 2.0
		{in: 2.0, se: 0x4000, m: 0x8000000000000000, acc: big.Exact},
		// 7FFE FFFFFFFF FFFFFFFF = 1.18973149535723176505e+4932 (max normal)
		//{in: 1.18973149535723176505e+4932, se: 0x7FFE, m: 0xFFFFFFFFFFFFFFFF, acc: big.Exact},
		// 0001 80000000 00000000 = 3.36210314311209350626e-4932 (min positive normal)
		//{in: 3.36210314311209350626e-4932, se: 0x0001, m: 0x8000000000000000, acc: big.Exact},
		// 0000 7FFFFFFF FFFFFFFF = 3.36210314311209350608e-4932 (max subnormal)
		//{in: 3.36210314311209350608e-4932, se: 0x0000, m: 0x7FFFFFFFFFFFFFFF, acc: big.Exact},
		// 0000 00000000 00000001 = 3.64519953188247460253e-4951 (min positive subnormal)
		//{in: 3.64519953188247460253e-4951, se: 0x0000, m: 0x0000000000000001, acc: big.Exact},
		// 7FFF 80000000 00000000 = infinity
		{in: math.Inf(1), se: 0x7FFF, m: 0x8000000000000000, acc: big.Exact},
		// FFFF 80000000 00000000 = -infinity
		{in: math.Inf(-1), se: 0xFFFF, m: 0x8000000000000000, acc: big.Exact},

		// 2^i
		// TODO: add test cases for 2^i
	}
	for _, g := range golden {
		f, acc := NewFromFloat64(g.in)
		se, m := f.Bits()
		if g.se != se || g.m != m {
			x, _ := f.Float64()
			t.Errorf("bits mismatch; expected 0x%04X %016X (%v), got 0x%04X %016X (%v)", g.se, g.m, g.in, se, m, x)
		}
		if g.acc != acc {
			x, _ := f.Float64()
			t.Errorf("accuracy mismatch; expected %v (%v), got %v (%v)", g.acc, g.in, acc, x)
		}
	}
}
