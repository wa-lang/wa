package bfloat

import (
	"fmt"
	"math/big"
)

const (
	// precision specifies the number of bits in the mantissa (including the
	// implicit lead bit).
	precision = 8
	// exponent bias.
	bias = 127
)

// Float is a floating-point number in bfloat16 floating-point format.
type Float struct {
	// Sign, exponent and fraction.
	//
	//    1 bit:   sign
	//    8 bits:  exponent
	//    7 bits:  fraction
	bits uint16
}

func NewFromBits(bits uint16) Float {
	return Float{bits: bits}
}

func (f Float) Big() (x *big.Float, nan bool) {
	signbit := f.Signbit()
	exp := f.Exp()
	frac := f.Frac()
	x = big.NewFloat(0)
	x.SetPrec(precision)
	x.SetMode(big.ToNearestEven)

	// ref: https://en.wikipedia.org/wiki/Bfloat16_floating-point_format#Contrast_with_bfloat16_and_single_precision
	//
	// 0b00001 - 0b11110
	// Normalized number.
	//
	//    (-1)^signbit * 2^(exp-127) * 1.mant_2
	lead := 1
	exponent := exp - bias

	switch exp {
	case 0xFF:
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
		//    (-1)^signbit * 2^(-126) * 0.mant_2
		lead = 0
		exponent = -126
	}

	// number = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
	sign := "+"
	if signbit {
		sign = "-"
	}
	s := fmt.Sprintf("%s0b%d.%07bp%d", sign, lead, frac, exponent)
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
	// 8 bit exponent: 0b0111111110000000
	return int(f.bits & 0x7F80 >> 7)
}

// Frac returns the fraction of f.
func (f Float) Frac() uint16 {
	// 7 bit mantissa: 0b0000000001111111
	return f.bits & 0x7F
}
