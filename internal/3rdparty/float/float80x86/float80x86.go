//go:generate go run gen.go -o extra_test.go

// Package float80x86 implements encoding and decoding of x86 extended precision
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Extended_precision#x86_extended_precision_format
package float80x86

import (
	"fmt"
	"math"
	"math/big"
)

const (
	// precision specifies the number of bits in the mantissa (including the
	// explicit lead bit).
	precision = 64
	// exponent bias.
	bias = 16383
)

// Positive and negative Not-a-Number, infinity and zero.
var (
	// +NaN
	NaN = Float{se: 0x7FFF, m: 0xBFFFFFFFFFFFFFFF}
	// -NaN
	NegNaN = Float{se: 0xFFFF, m: 0xBFFFFFFFFFFFFFFF}
	// +Inf
	Inf = Float{se: 0x7FFF, m: 0x8000000000000000}
	// -Inf
	NegInf = Float{se: 0xFFFF, m: 0x8000000000000000}
	// +zero
	Zero = Float{se: 0x0000, m: 0x0000000000000000}
	// -zero
	NegZero = Float{se: 0x8000, m: 0x0000000000000000}
)

// Float is a floating-point number in x86 extended precision format.
type Float struct {
	// Sign and exponent.
	//
	//    1 bit:   sign
	//    15 bits: exponent
	se uint16
	// Integer part and fraction.
	//
	//    1 bit:   integer part
	//    63 bits: fraction
	m uint64
}

// NewFromBits returns the floating-point number corresponding to the x86
// extended precision representation.
func NewFromBits(se uint16, m uint64) Float {
	return Float{se: se, m: m}
}

// NewFromFloat32 returns the nearest x86 extended precision floating-point
// number for x and the accuracy of the conversion.
func NewFromFloat32(x float32) (Float, big.Accuracy) {
	f, acc := NewFromFloat64(float64(x))
	if acc == big.Exact {
		_, acc = f.Float32()
	}
	return f, acc
}

// NewFromFloat64 returns the nearest x86 extended precision floating-point
// number for x and the accuracy of the conversion.
func NewFromFloat64(x float64) (Float, big.Accuracy) {
	// +-NaN
	switch {
	case math.IsNaN(x):
		if math.Signbit(x) {
			// -NaN
			//    sign: 1
			//    exp:  all ones
			//    mant: 10 non-zero
			return NegNaN, big.Exact
		}
		// +NaN
		//    sign: 0
		//    exp:  all ones
		//    mant: 10 non-zero
		return NaN, big.Exact
	}
	y := big.NewFloat(x)
	y.SetPrec(precision)
	y.SetMode(big.ToNearestEven)
	// TODO: check accuracy after setting precision?
	return NewFromBig(y)
}

// NewFromBig returns the nearest x86 extended precision floating-point number
// for x and the accuracy of the conversion.
func NewFromBig(x *big.Float) (Float, big.Accuracy) {
	// +-Inf
	zero := big.NewFloat(0)
	switch {
	case x.IsInf():
		if x.Signbit() {
			// -Inf
			//    sign: 1
			//    exp:  all ones
			//    mant: 10 zero
			return NegInf, big.Exact
		}
		// +Inf
		//    sign: 0
		//    exp:  all ones
		//    mant: 10 zero
		return Inf, big.Exact
	// +-zero
	case x.Cmp(zero) == 0:
		if x.Signbit() {
			// -zero
			//    sign: 1
			//    exp:  zero
			//    mant: zero
			return NegZero, big.Exact
		}
		// +zero
		//    sign: 0
		//    exp:  zero
		//    mant: zero
		return Zero, big.Exact
	}

	// Sign
	var se uint16
	if x.Signbit() {
		se |= 0x8000
	}

	// Exponent and mantissa.
	var m uint64
	mant := &big.Float{}
	exponent := x.MantExp(mant)
	// TODO: verify, as float80x86 also has an explicit lead bit.
	// Remove 1 from the exponent as big.Float has an no lead bit.
	exp := exponent - 1 + bias

	// Handle denormalized values.
	// TODO: validate implementation of denormalized values.
	if exp <= 0 {
		acc := big.Exact
		if exp <= -(precision - 1) {
			exp = precision - 1
			acc = big.Below
		}
		mant.SetMantExp(mant, exp+precision-1)
		if mant.Signbit() {
			mant.Neg(mant)
		}
		v, _ := mant.Uint64()
		// TODO: calculate acc based on if v&^0x7FFFFFFFFFFFFFFF != 0 {}
		m |= v & 0x7FFFFFFFFFFFFFFF
		return Float{se: se, m: m}, acc
	}

	// 0b111111111111111
	acc := big.Exact
	if (exp &^ 0x7FFF) != 0 {
		acc = big.Above
	}
	se |= uint16(exp & 0x7FFF)

	if mant.Signbit() {
		mant.Neg(mant)
	}
	mant.SetMantExp(mant, precision)
	if !mant.IsInt() {
		acc = big.Below
	}
	// mantissa, including explicit lead bit
	mantissa, acc2 := mant.Uint64()
	if acc == big.Exact {
		acc = acc2
	}
	m |= mantissa
	return Float{se: se, m: m}, acc
}

// Bits returns the x86 extended precision binary representation of f.
func (f Float) Bits() (se uint16, m uint64) {
	return f.se, f.m
}

// Float32 returns the float32 value nearest to f. If f is too small to be
// represented by a float32 (|f| < math.SmallestNonzeroFloat32), the result is
// (0, Below) or (-0, Above), respectively, depending on the sign of f. If f is
// too large to be represented by a float32 (|f| > math.MaxFloat32), the result
// is (+Inf, Above) or (-Inf, Below), depending on the sign of f.
func (f Float) Float32() (float32, big.Accuracy) {
	x, nan := f.Big()
	if nan {
		if x.Signbit() {
			return float32(-math.NaN()), big.Exact
		}
		return float32(math.NaN()), big.Exact
	}
	return x.Float32()
}

// Float64 returns the float64 value nearest to f. If f is too small to be
// represented by a float64 (|f| < math.SmallestNonzeroFloat64), the result is
// (0, Below) or (-0, Above), respectively, depending on the sign of f. If f is
// too large to be represented by a float64 (|f| > math.MaxFloat64), the result
// is (+Inf, Above) or (-Inf, Below), depending on the sign of f.
func (f Float) Float64() (float64, big.Accuracy) {
	x, nan := f.Big()
	if nan {
		if x.Signbit() {
			return -math.NaN(), big.Exact
		}
		return math.NaN(), big.Exact
	}
	return x.Float64()
}

// Big returns the multi-precision floating-point number representation of f and
// a boolean indicating whether f is Not-a-Number.
func (f Float) Big() (x *big.Float, nan bool) {
	signbit := f.Signbit()
	exp := f.Exp()
	x = big.NewFloat(0)
	x.SetPrec(precision)
	x.SetMode(big.ToNearestEven)

	// ref: https://en.wikipedia.org/wiki/Extended_precision#x86_extended_precision_format
	//
	// 0b000000000000001 - 0b111111111111110
	// Normalized number.
	//
	//    (-1)^signbit * 2^(exp-16383) * 1.mant_2
	exponent := exp - bias

	switch exp {
	// 0b111111111111111
	case 0x7FFF:
		// Inf or NaN
		if f.m == 0x8000000000000000 {
			// +-Inf
			//    10 zero
			x.SetInf(signbit)
			return x, false
		}
		// +-NaN
		//    10 non-zero
		if signbit {
			x.Neg(x)
		}
		return x, true
	// 0b000000000000000
	case 0x0000:
		if f.m == 0 {
			// +-Zero
			if signbit {
				x.Neg(x)
			}
			return x, false
		}
		// Denormalized number.
		//
		//    (-1)^signbit * 2^(-16382) * 0.mant_2
		exponent = -16382
	}

	// number = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
	sign := "+"
	if signbit {
		sign = "-"
	}
	lead := f.Lead()
	frac := f.Frac()
	s := fmt.Sprintf("%s0b%d.%063bp%d", sign, lead, frac, exponent)
	if _, _, err := x.Parse(s, 0); err != nil {
		panic(err)
	}
	return x, false
}

// Signbit reports whether f is negative or negative 0.
func (f Float) Signbit() bool {
	// 0b1000000000000000
	return f.se&0x8000 != 0
}

// Exp returns the exponent of f.
func (f Float) Exp() int {
	// 0b0111111111111111
	return int(f.se & 0x7FFF)
}

// Lead returns the explicit lead bit of f.
func (f Float) Lead() int {
	return int(f.m >> 63)
}

// Frac returns the fraction of f.
func (f Float) Frac() uint64 {
	return f.m & 0x7FFFFFFFFFFFFFFF
}
