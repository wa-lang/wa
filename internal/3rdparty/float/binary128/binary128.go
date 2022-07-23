//go:generate go run gen.go -o extra_test.go

// Package binary128 implements encoding and decoding of IEEE 754 quadruple
// precision floating-point numbers.
//
// https://en.wikipedia.org/wiki/Quadruple-precision_floating-point_format
package binary128

import (
	"fmt"
	"math"
	"math/big"
)

const (
	// precision specifies the number of bits in the mantissa (including the
	// implicit lead bit).
	precision = 113
	// exponent bias.
	bias = 16383
)

// Positive and negative Not-a-Number, infinity and zero.
var (
	// +NaN
	NaN = Float{a: 0x7FFF800000000000, b: 0}
	// -NaN
	NegNaN = Float{a: 0xFFFF800000000000, b: 0}
	// +Inf
	Inf = Float{a: 0x7FFF000000000000, b: 0}
	// -Inf
	NegInf = Float{a: 0xFFFF000000000000, b: 0}
	// +zero
	Zero = Float{a: 0x0000000000000000, b: 0}
	// -zero
	NegZero = Float{a: 0x8000000000000000, b: 0}
)

// Float is a floating-point number in IEEE 754 quadruple precision format.
type Float struct {
	// Sign, exponent and fraction.
	//
	//    1 bit:    sign
	//    15 bits:  exponent
	//    112 bits: fraction
	a uint64
	b uint64
}

// NewFromBits returns the floating-point number corresponding to the IEEE 754
// quadruple precision binary representation.
func NewFromBits(a, b uint64) Float {
	return Float{a: a, b: b}
}

// NewFromFloat32 returns the nearest quadruple precision floating-point number
// for x and the accuracy of the conversion.
func NewFromFloat32(x float32) (Float, big.Accuracy) {
	f, acc := NewFromFloat64(float64(x))
	if acc == big.Exact {
		_, acc = f.Float32()
	}
	return f, acc
}

// NewFromFloat64 returns the nearest quadruple precision floating-point number
// for x and the accuracy of the conversion.
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

// NewFromBig returns the nearest quadruple precision floating-point number for
// x and the accuracy of the conversion.
func NewFromBig(x *big.Float) (Float, big.Accuracy) {
	// +-Inf
	zero := big.NewFloat(0).SetPrec(precision)
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
	var a, b uint64
	if x.Signbit() {
		a |= 0x8000000000000000
	}

	// Exponent and mantissa.
	mant := new(big.Float).SetPrec(precision)
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
		mantissa, _ := mant.Int(nil)
		maskA := big.NewInt(0)
		for i := 64; i < 112; i++ {
			maskA.SetBit(maskA, i, 1)
		}
		maskB := big.NewInt(0)
		for i := 0; i < 64; i++ {
			maskB.SetBit(maskB, i, 1)
		}
		bigA := new(big.Int).And(mantissa, maskA) // a = (mantissa & maskA) >> 64
		bigA.Rsh(bigA, 64)
		bigB := new(big.Int).And(mantissa, maskB) // b = mantissa & maskB
		// TODO: calculate acc based on if mantissa&^maskA != 0 {}
		a |= bigA.Uint64() & 0x0000FFFFFFFFFFFF
		b = bigB.Uint64()
		return Float{a: a, b: b}, acc
	}

	// exponent mask (15 bits): 0b111111111111111
	acc := big.Exact
	if (exp &^ 0x7FFF) != 0 {
		acc = big.Above
	}
	a |= uint64(exp&0x7FFF) << 48

	if mant.Signbit() {
		mant.Neg(mant)
	}
	mant.SetMantExp(mant, precision)
	if !mant.IsInt() {
		acc = big.Below
	}
	mantissa, _ := mant.Int(nil)
	mantissa.SetBit(mantissa, 112, 0) // clear implicit lead bit; 2^112

	// mantissa mask (113 bits, including implicit lead bit): 0x1FFFFFFFFFFFFFFFFFFFFFFFFFFFF
	maskA := big.NewInt(0)
	for i := 64; i < 112; i++ {
		maskA.SetBit(maskA, i, 1)
	}
	maskB := big.NewInt(0)
	for i := 0; i < 64; i++ {
		maskB.SetBit(maskB, i, 1)
	}
	bigA := new(big.Int).And(mantissa, maskA) // a = (mantissa & maskA) >> 64
	bigA.Rsh(bigA, 64)
	bigB := new(big.Int).And(mantissa, maskB) // b = mantissa & maskB
	if acc == big.Exact && (bigA.Uint64()&^0x0000FFFFFFFFFFFF) != 0 {
		acc = big.Below
	}
	a |= bigA.Uint64() & 0x0000FFFFFFFFFFFF
	b = bigB.Uint64()
	return Float{a: a, b: b}, acc
}

// Bits returns the IEEE 754 quadruple precision binary representation of f.
func (f Float) Bits() (a, b uint64) {
	return f.a, f.b
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
	frac1, frac2 := f.Frac()
	x = big.NewFloat(0)
	x.SetPrec(precision)
	x.SetMode(big.ToNearestEven)

	lead := 1
	exponent := exp - bias

	switch exp {
	// 0b111111111111111
	case 0x7FFF:
		// Inf or NaN
		if frac1 == 0 && frac2 == 0 {
			// +-Inf
			x.SetInf(signbit)
			return x, false
		}
		// +-NaN
		if signbit {
			x.Neg(x)
		}
		return x, true
	// 0b000000000000000
	case 0x0000:
		if frac1 == 0 && frac2 == 0 {
			// +-Zero
			if signbit {
				x.Neg(x)
			}
			return x, false
		}
		// Denormalized number.
		//
		//    (-1)^signbit * 2^(-16382) * 0.mant_2
		lead = 0
		exponent = -16382
	}

	// number = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
	sign := "+"
	if signbit {
		sign = "-"
	}
	// first part cut the sign and exponent which only contains 48 bits
	fracStr := fmt.Sprintf("%048b%064b", frac1, frac2)
	s := fmt.Sprintf("%s0b%d.%sp%d", sign, lead, fracStr, exponent)
	if _, _, err := x.Parse(s, 0); err != nil {
		panic(err)
	}
	return x, false
}

// Signbit reports whether f is negative or negative 0.
func (f Float) Signbit() bool {
	// first bit is sign bit
	return f.a&0x8000000000000000 != 0
}

// Exp returns the exponent of f.
func (f Float) Exp() int {
	// 15 bit exponent
	return int(f.a&0x7FFF000000000000) >> 48
}

// Frac returns the fraction of f.
func (f Float) Frac() (uint64, uint64) {
	// 0x0000FFFFFFFFFFFF removes the sign and exponent part (total 16 bits) from
	// our floating-point number. Now we can say it contains 48 bits of fraction,
	// and `f.b` part has the rest of fraction.
	return (f.a & 0x0000FFFFFFFFFFFF), f.b
}
