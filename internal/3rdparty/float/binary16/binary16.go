//go:generate go run gen.go -o extra_test.go

// Package binary16 implements encoding and decoding of IEEE 754 half precision
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Half-precision_floating-point_format
package binary16

import (
	"fmt"
	"math"
	"math/big"
)

const (
	// precision specifies the number of bits in the mantissa (including the
	// implicit lead bit).
	precision = 11
	// exponent bias.
	bias = 15
)

// Positive and negative Not-a-Number, infinity and zero.
var (
	// +NaN
	NaN = Float{bits: 0x7E00}
	// -NaN
	NegNaN = Float{bits: 0xFE00}
	// +Inf
	Inf = Float{bits: 0x7C00}
	// -Inf
	NegInf = Float{bits: 0xFC00}
	// +zero
	Zero = Float{bits: 0x0000}
	// -zero
	NegZero = Float{bits: 0x8000}
)

// Float is a floating-point number in IEEE 754 half precision format.
type Float struct {
	// Sign, exponent and fraction.
	//
	//    1 bit:   sign
	//    5 bits:  exponent
	//    10 bits: fraction
	bits uint16
}

// NewFromBits returns the floating-point number corresponding to the IEEE 754
// half precision binary representation.
func NewFromBits(bits uint16) Float {
	return Float{bits: bits}
}

// NewFromFloat32 returns the nearest half precision floating-point number for x
// and the accuracy of the conversion.
func NewFromFloat32(x float32) (Float, big.Accuracy) {
	f, acc := NewFromFloat64(float64(x))
	if acc == big.Exact {
		_, acc = f.Float32()
	}
	return f, acc
}

// NewFromFloat64 returns the nearest half precision floating-point number for x
// and the accuracy of the conversion.
func NewFromFloat64(x float64) (Float, big.Accuracy) {
	// +-NaN
	switch {
	case math.IsNaN(x):
		if math.Signbit(x) {
			// -NaN
			return NegNaN, big.Exact
		}
		// +NaN
		return NaN, big.Exact
	}
	y := big.NewFloat(x)
	y.SetPrec(precision)
	y.SetMode(big.ToNearestEven)
	// TODO: check accuracy after setting precision?
	return NewFromBig(y)
}

// NewFromBig returns the nearest half precision floating-point number for x and
// the accuracy of the conversion.
func NewFromBig(x *big.Float) (Float, big.Accuracy) {
	// +-Inf
	zero := big.NewFloat(0)
	switch {
	case x.IsInf():
		if x.Signbit() {
			// -Inf
			return NegInf, big.Exact
		}
		// +Inf
		return Inf, big.Exact
	// +-zero
	case x.Cmp(zero) == 0:
		if x.Signbit() {
			// -zero
			return NegZero, big.Exact
		}
		// +zero
		return Zero, big.Exact
	}

	// Sign
	var bits uint16
	if x.Signbit() {
		bits |= 0x8000
	}

	// Exponent and mantissa.
	mant := new(big.Float)
	exponent := x.MantExp(mant)
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
		mantissa, _ := mant.Uint64()
		// TODO: calculate acc based on if mantissa&^0x7FF != 0 {}
		bits |= uint16(mantissa & 0x7FF)
		return Float{bits: bits}, acc
	}

	// exponent mask (5 bits): 0b11111
	acc := big.Exact
	if (exp &^ 0x1F) != 0 {
		acc = big.Above
	}
	bits |= uint16(exp&0x1F) << 10

	if mant.Signbit() {
		mant.Neg(mant)
	}
	mant.SetMantExp(mant, precision)
	if !mant.IsInt() {
		acc = big.Below
	}
	mantissa, _ := mant.Uint64()
	mantissa &^= 0x400 // clear implicit lead bit; 2^10

	// mantissa mask (11 bits, including implicit lead bit): 0b11111111111
	if acc == big.Exact && (mantissa&^0x7FF) != 0 {
		acc = big.Below
	}
	mantissa &= 0x7FF
	bits |= uint16(mantissa)
	return Float{bits: bits}, acc
}

// Bits returns the IEEE 754 half precision binary representation of f.
func (f Float) Bits() uint16 {
	return f.bits
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
	frac := f.Frac()
	x = big.NewFloat(0)
	x.SetPrec(precision)
	x.SetMode(big.ToNearestEven)

	// ref: https://en.wikipedia.org/wiki/Half-precision_floating-point_format#Exponent_encoding
	//
	// 0b00001 - 0b11110
	// Normalized number.
	//
	//    (-1)^signbit * 2^(exp-15) * 1.mant_2
	lead := 1
	exponent := exp - bias

	switch exp {
	// 0b11111
	case 0x1F:
		// Inf or NaN
		if frac == 0 {
			// +-Inf
			x.SetInf(signbit)
			return x, false
		}
		// +-NaN
		if signbit {
			x.Neg(x)
		}
		return x, true
	// 0b00000
	case 0x00:
		if frac == 0 {
			// +-Zero
			if signbit {
				x.Neg(x)
			}
			return x, false
		}
		// Denormalized number.
		//
		//    (-1)^signbit * 2^(-14) * 0.mant_2
		lead = 0
		exponent = -14
	}

	// number = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
	sign := "+"
	if signbit {
		sign = "-"
	}
	s := fmt.Sprintf("%s0b%d.%010bp%d", sign, lead, frac, exponent)
	if _, _, err := x.Parse(s, 0); err != nil {
		panic(err)
	}
	return x, false
}

// Signbit reports whether f is negative or negative 0.
func (f Float) Signbit() bool {
	// first bit is sign bit: 0b1000000000000000
	return f.bits&0x8000 != 0
}

// Exp returns the exponent of f.
func (f Float) Exp() int {
	// 5 bit exponent: 0b0111110000000000
	return int(f.bits & 0x7C00 >> 10)
}

// Frac returns the fraction of f.
func (f Float) Frac() uint16 {
	// 10 bit mantissa: 0b0000001111111111
	return f.bits & 0x03FF
}
