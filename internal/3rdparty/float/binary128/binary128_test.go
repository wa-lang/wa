package binary128

import (
	"fmt"
	"math"
	"math/big"
	"testing"
)

func TestNewFromBits(t *testing.T) {
	const rawpi = "3.1415926535897932384626433832795028"
	pi, ok := newFloat(0).SetString(rawpi)
	if !ok {
		panic(fmt.Errorf("unable to create arbitrary floating-point value of pi (%q)", rawpi))
	}
	golden := []struct {
		a, b uint64
		want *big.Float
		nan  bool
	}{
		// Special numbers.
		// +NaN
		// 0x7FFF 8000000000000000000000000000 = +NaN
		{a: 0x7FFF800000000000, b: 0x0000000000000001, want: newFloat(0), nan: true},
		// -NaN
		// 0xFFFF 8000000000000000000000000000 = -NaN
		{a: 0xFFFF800000000000, b: 0x0000000000000000, want: newFloat(math.Copysign(0, -1)), nan: true},
		// +inf
		// 0x7FFF0000000000000000000000000000 = +inf
		{a: 0x7FFF000000000000, b: 0x0000000000000000, want: newFloat(0).SetInf(false)},
		// -inf
		// 0xFFFF0000000000000000000000000000 = -inf
		{a: 0xFFFF000000000000, b: 0x0000000000000000, want: newFloat(0).SetInf(true)},
		// +0
		// 0x00000000000000000000000000000000 = +0
		{a: 0x0000000000000000, b: 0x0000000000000000, want: newFloat(+0)},
		// -0
		// 0x80000000000000000000000000000000 = -0
		{a: 0x8000000000000000, b: 0x0000000000000000, want: newFloat(math.Copysign(0, -1))},

		// from: https://en.wikipedia.org/wiki/Quadruple-precision_floating-point_format#Quadruple_precision_examples

		// smallest positive subnormal number
		// 0x00000000000000000000000000000001 = 2^{-16382} * 2^{-112} =  2^{-16494}
		{a: 0x0000000000000000, b: 0x0000000000000001, want: pow(newFloat(2), -16494)},
		// largest subnormal number
		// 0x0000FFFFFFFFFFFFFFFFFFFFFFFFFFFF = 2^{-16382} * (1 - 2^{-112})
		{a: 0x0000FFFFFFFFFFFF, b: 0xFFFFFFFFFFFFFFFF, want: mul(pow(newFloat(2), -16382), sub(newFloat(1), pow(newFloat(2), -112)))},
		// smallest positive normal number
		// 0x00010000000000000000000000000000 = 2^{-16382}
		{a: 0x0001000000000000, b: 0x0000000000000000, want: pow(newFloat(2), -16382)},
		// largest normal number
		// 0x7FFEFFFFFFFFFFFFFFFFFFFFFFFFFFFF = 2^16383 * (2 - 2^{-112})
		{a: 0x7FFEFFFFFFFFFFFF, b: 0xFFFFFFFFFFFFFFFF, want: mul(pow(newFloat(2), 16383), sub(newFloat(2), pow(newFloat(2), -112)))},
		// largest number less than one
		// 0x3FFEFFFFFFFFFFFFFFFFFFFFFFFFFFFF = 1 - 2^{-113}
		{a: 0x3FFEFFFFFFFFFFFF, b: 0xFFFFFFFFFFFFFFFF, want: sub(newFloat(1), pow(newFloat(2), -113))},
		// one
		// 0x3FFF0000000000000000000000000000 = 1
		{a: 0x3FFF000000000000, b: 0x0000000000000000, want: newFloat(1)},
		// smallest number larger than one
		// 0x3FFF0000000000000000000000000001 = 1 + 2^{-112}
		{a: 0x3FFF000000000000, b: 0x0000000000000001, want: add(newFloat(1), pow(newFloat(2), -112))},
		// -2
		// 0xC0000000000000000000000000000000 = -2
		{a: 0xC000000000000000, b: 0x0000000000000000, want: newFloat(-2)},
		// pi
		// 0x4000921FB54442D18469898CC51701B8 = pi
		{a: 0x4000921FB54442D1, b: 0x8469898CC51701B8, want: pi},
		// 1/3
		// 0x3FFD5555555555555555555555555555 = 1/3
		{a: 0x3FFD555555555555, b: 0x5555555555555555, want: newFloat(0).SetRat(big.NewRat(1, 3))},
	}
	for _, g := range golden {
		f := NewFromBits(g.a, g.b)
		got, nan := f.Big()
		if g.want.Cmp(got) != 0 {
			t.Errorf("0x%016X%016X: floating-point number mismatch; expected %v, got %v", g.a, g.b, g.want, got)
		}
		if g.nan != nan {
			t.Errorf("0x%016X%016X: floating-point Not-a-Number indicator mismatch; expected %v, got %v", g.a, g.b, g.nan, nan)
		}
	}
}

func TestNewFromFloat64(t *testing.T) {
	golden := []struct {
		in   float64
		a, b uint64
		acc  big.Accuracy
	}{
		// Special numbers.
		// 0x7FFF 8000000000000000000000000000 = +NaN
		{in: math.NaN(), a: 0x7FFF800000000000, b: 0x0000000000000000, acc: big.Exact},
		// -NaN
		// 0xFFFF 8000000000000000000000000000 = -NaN
		{in: -math.NaN(), a: 0xFFFF800000000000, b: 0x0000000000000000, acc: big.Exact},
		// +inf
		// 0x7FFF0000000000000000000000000000 = +inf
		{in: math.Inf(+1), a: 0x7FFF000000000000, b: 0x0000000000000000, acc: big.Exact},
		// -inf
		// 0xFFFF0000000000000000000000000000 = -inf
		{in: math.Inf(-1), a: 0xFFFF000000000000, b: 0x0000000000000000, acc: big.Exact},
		// +0
		// 0x00000000000000000000000000000000 = +0
		{in: +0, a: 0x0000000000000000, b: 0x0000000000000000, acc: big.Exact},
		// -0
		// 0x80000000000000000000000000000000 = -0
		{in: math.Copysign(0, -1), a: 0x8000000000000000, b: 0x0000000000000000, acc: big.Exact},

		// from: https://en.wikipedia.org/wiki/Quadruple-precision_floating-point_format#Quadruple_precision_examples

		// one
		// 0x3FFF0000000000000000000000000000 = 1
		{in: 1, a: 0x3FFF000000000000, b: 0x0000000000000000, acc: big.Exact},
		// -2
		// 0xC0000000000000000000000000000000 = -2
		{in: -2, a: 0xC000000000000000, b: 0x0000000000000000, acc: big.Exact},
		// pi
		// 0x4000921FB54442D18469898CC51701B8 = pi
		{in: math.Pi, a: 0x4000921FB54442D1, b: 0x8469898CC51701B8, acc: big.Exact},
		// 1/3
		// 0x3FFD5555555555555555555555555555 = 1/3
		{in: 1.0 / 3.0, a: 0x3FFD555555555555, b: 0x5555555555555555, acc: big.Exact},
	}
	for _, g := range golden {
		f, acc := NewFromFloat64(g.in)
		a, b := f.Bits()
		x, _ := f.Float64()
		// mask last 60 bits, as float64 only has 53 bits precision (as compared to
		// 113 of binary128).
		const mask = 0xF000000000000000
		wantBMask := g.b & mask
		gotBMask := b & mask
		if g.a != a || wantBMask != gotBMask {
			t.Errorf("bits mismatch; expected 0x%016X%016X (%v), got 0x%016X%016X (%v)", g.a, wantBMask, g.in, a, gotBMask, x)
		}
		if g.acc != acc {
			t.Errorf("accuracy mismatch; expected %v (%v), got %v (%v)", g.acc, g.in, acc, x)
		}
	}
}

// ### [ Helper functions ] ####################################################

// pow returns x**y, the base-x exponential of y.
func pow(x *big.Float, y int64) *big.Float {
	switch {
	// x^{-42}
	case y < 0:
		z := newFloat(1)
		for i := int64(0); i < -y; i++ {
			z = div(z, x)
		}
		return z
	// x^42
	case y > 0:
		z := newFloat(1)
		for i := int64(0); i < y; i++ {
			z = mul(z, x)
		}
		return z
	// x^0
	default: // y == 0
		return newFloat(1)
	}
}

// add returns the sum x+y.
func add(x, y *big.Float) *big.Float {
	return newFloat(0).Add(x, y)
}

// add returns the difference x-y.
func sub(x, y *big.Float) *big.Float {
	return newFloat(0).Sub(x, y)
}

// add returns the product x*y.
func mul(x, y *big.Float) *big.Float {
	return newFloat(0).Mul(x, y)
}

// add returns the quotient x/y.
func div(x, y *big.Float) *big.Float {
	return newFloat(0).Quo(x, y)
}

// newFloat returns a new floating-point value based on x with precision 113.
func newFloat(x float64) *big.Float {
	return big.NewFloat(0).SetPrec(precision).SetFloat64(x)
}
