package binary16

import (
	"math"
	"math/big"
	"testing"
)

func TestNewFromBits(t *testing.T) {
	golden := []struct {
		bits uint16
		want float64
	}{
		// Special numbers.
		// 0 11111 1000000000 = +NaN
		{bits: 0x7E00, want: math.NaN()},
		// -NaN
		// 1 11111 1000000000 = -NaN
		{bits: 0xFE00, want: -math.NaN()},

		// from: https://en.wikipedia.org/wiki/Half-precision_floating-point_format#Half_precision_examples

		// 0 01111 0000000000 = 1
		{bits: 0x3C00, want: 1},
		// 0 01111 0000000001 = 1 + 2^(-10) = 1.0009765625 (next smallest float after 1)
		{bits: 0x3C01, want: 1.0009765625},
		// 1 10000 0000000000 = -2
		{bits: 0xC000, want: -2},
		// 0 11110 1111111111 = 65504 (max half precision)
		{bits: 0x7BFF, want: 65504},
		// 0 00001 0000000000 = 2^(-14) ~= 6.10352 * 10^(-5) (minimum positive normal)
		{bits: 0x0400, want: math.Pow(2, -14)},
		// 0 00000 0000000001 = 2^(-24) ~= 5.96046 * 10^(-8) (minimum positive subnormal)
		{bits: 0x0001, want: math.Pow(2, -24)},
		// 0 00000 0000000000 = 0
		{bits: 0x0000, want: 0},
		// 1 00000 0000000000 = −0
		{bits: 0x8000, want: math.Copysign(0, -1)},
		// 0 11111 0000000000 = infinity
		{bits: 0x7C00, want: math.Inf(1)},
		// 1 11111 0000000000 = -infinity
		{bits: 0xFC00, want: math.Inf(-1)},
		// 0 01101 0101010101 = 0.333251953125 ~= 1/3
		{bits: 0x3555, want: 0.333251953125},

		// from: https://reviews.llvm.org/rL237161

		// Normalized numbers.
		// 0 01110 0000000000 = 0.5
		{bits: 0x3800, want: 0.5},
		// 1 01110 0000000000 = -0.5
		{bits: 0xB800, want: -0.5},
		// 0 01111 1000000000 = 1.5
		{bits: 0x3E00, want: 1.5},
		// 1 01111 1000000000 = -1.5
		{bits: 0xBE00, want: -1.5},
		// 0 10000 0100000000 = 2.5
		{bits: 0x4100, want: 2.5},
		// 1 10000 0100000000 = -2.5
		{bits: 0xC100, want: -2.5},
		// Denormalized numbers.
		// 0 00000 0000010000 = 2^(-20)
		{bits: 0x0010, want: math.Pow(2, -20)},
		// 1 00000 0000000001 = -2^(-24)
		{bits: 0x8001, want: -math.Pow(2, -24)},

		// 2^i
		{bits: 0x0001, want: math.Pow(2, -24)}, // 2^(-24)
		{bits: 0x0002, want: math.Pow(2, -23)}, // 2^(-23)
		{bits: 0x0004, want: math.Pow(2, -22)}, // 2^(-22)
		{bits: 0x0008, want: math.Pow(2, -21)}, // 2^(-21)
		{bits: 0x0010, want: math.Pow(2, -20)}, // 2^(-20)
		{bits: 0x0020, want: math.Pow(2, -19)}, // 2^(-19)
		{bits: 0x0040, want: math.Pow(2, -18)}, // 2^(-18)
		{bits: 0x0080, want: math.Pow(2, -17)}, // 2^(-17)
		{bits: 0x0100, want: math.Pow(2, -16)}, // 2^(-16)
		{bits: 0x0200, want: math.Pow(2, -15)}, // 2^(-15)
		{bits: 0x0400, want: math.Pow(2, -14)}, // 2^(-14)
		{bits: 0x0800, want: math.Pow(2, -13)}, // 2^(-13)
		{bits: 0x0C00, want: math.Pow(2, -12)}, // 2^(-12)
		{bits: 0x1000, want: math.Pow(2, -11)}, // 2^(-11)
		{bits: 0x1400, want: math.Pow(2, -10)}, // 2^(-10)
		{bits: 0x1800, want: math.Pow(2, -9)},  // 2^(-9)
		{bits: 0x1C00, want: math.Pow(2, -8)},  // 2^(-8)
		{bits: 0x2000, want: math.Pow(2, -7)},  // 2^(-7)
		{bits: 0x2400, want: math.Pow(2, -6)},  // 2^(-6)
		{bits: 0x2800, want: math.Pow(2, -5)},  // 2^(-5)
		{bits: 0x2C00, want: math.Pow(2, -4)},  // 2^(-4)
		{bits: 0x3000, want: math.Pow(2, -3)},  // 2^(-3)
		{bits: 0x3400, want: math.Pow(2, -2)},  // 2^(-2)
		{bits: 0x3800, want: math.Pow(2, -1)},  // 2^(-1)
		{bits: 0x3C00, want: math.Pow(2, 0)},   // 2^0
		{bits: 0x4000, want: math.Pow(2, 1)},   // 2^1
		{bits: 0x4400, want: math.Pow(2, 2)},   // 2^2
		{bits: 0x4800, want: math.Pow(2, 3)},   // 2^3
		{bits: 0x4C00, want: math.Pow(2, 4)},   // 2^4
		{bits: 0x5000, want: math.Pow(2, 5)},   // 2^5
		{bits: 0x5400, want: math.Pow(2, 6)},   // 2^6
		{bits: 0x5800, want: math.Pow(2, 7)},   // 2^7
		{bits: 0x5C00, want: math.Pow(2, 8)},   // 2^8
		{bits: 0x6000, want: math.Pow(2, 9)},   // 2^9
		{bits: 0x6400, want: math.Pow(2, 10)},  // 2^10
		{bits: 0x6800, want: math.Pow(2, 11)},  // 2^11
		{bits: 0x6C00, want: math.Pow(2, 12)},  // 2^12
		{bits: 0x7000, want: math.Pow(2, 13)},  // 2^13
		{bits: 0x7400, want: math.Pow(2, 14)},  // 2^14
		{bits: 0x7800, want: math.Pow(2, 15)},  // 2^15
	}
	for _, g := range golden {
		f := NewFromBits(g.bits)
		got, _ := f.Float64()
		wantBits := math.Float64bits(g.want)
		gotBits := math.Float64bits(got)
		//fmt.Printf("bits: 0x%04X (%v)\n", g.bits, g.want)
		if wantBits != gotBits {
			t.Errorf("0x%04X: number mismatch; expected 0x%016X (%v), got 0x%016X (%v)", g.bits, wantBits, g.want, gotBits, got)
		}
	}
}

func TestNewFromFloat64(t *testing.T) {
	golden := []struct {
		in   float64
		want uint16
		acc  big.Accuracy
	}{
		// Special numbers.
		// 0 11111 1000000000 = +NaN
		{in: math.NaN(), want: 0x7E00, acc: big.Exact},
		// -NaN
		// 1 11111 1000000000 = -NaN
		{in: -math.NaN(), want: 0xFE00, acc: big.Exact},

		// from: https://en.wikipedia.org/wiki/Half-precision_floating-point_format#Half_precision_examples

		// 0 01111 0000000000 = 1
		{in: 1, want: 0x3C00, acc: big.Exact},
		// 0 01111 0000000001 = 1 + 2^(-10) = 1.0009765625 (next smallest float after 1)
		{in: 1.0009765625, want: 0x3C01, acc: big.Exact},
		// 1 10000 0000000000 = -2
		{in: -2, want: 0xC000, acc: big.Exact},
		// 0 11110 1111111111 = 65504 (max half precision)
		{in: 65504, want: 0x7BFF, acc: big.Exact},
		// 0 00001 0000000000 = 2^(-14) ~= 6.10352 * 10^(-5) (minimum positive normal)
		{in: math.Pow(2, -14), want: 0x0400, acc: big.Exact},
		// 0 00000 0000000001 = 2^(-24) ~= 5.96046 * 10^(-8) (minimum positive subnormal)
		{in: math.Pow(2, -24), want: 0x0001, acc: big.Exact},
		// 0 00000 0000000000 = 0
		{in: 0, want: 0x0000, acc: big.Exact},
		// 1 00000 0000000000 = −0
		{in: math.Copysign(0, -1), want: 0x8000, acc: big.Exact},
		// 0 11111 0000000000 = infinity
		{in: math.Inf(1), want: 0x7C00, acc: big.Exact},
		// 1 11111 0000000000 = -infinity
		{in: math.Inf(-1), want: 0xFC00, acc: big.Exact},
		// 0 01101 0101010101 = 0.333251953125 ~= 1/3
		{in: 0.333251953125, want: 0x3555, acc: big.Exact},

		// from: https://reviews.llvm.org/rL237161

		// Normalized numbers.
		// 0 01110 0000000000 = 0.5
		{in: 0.5, want: 0x3800, acc: big.Exact},
		// 1 01110 0000000000 = -0.5
		{in: -0.5, want: 0xB800, acc: big.Exact},
		// 0 01111 1000000000 = 1.5
		{in: 1.5, want: 0x3E00, acc: big.Exact},
		// 1 01111 1000000000 = -1.5
		{in: -1.5, want: 0xBE00, acc: big.Exact},
		// 0 10000 0100000000 = 2.5
		{in: 2.5, want: 0x4100, acc: big.Exact},
		// 1 10000 0100000000 = -2.5
		{in: -2.5, want: 0xC100, acc: big.Exact},
		// Denormalized numbers.
		// 0 00000 0000010000 = 2^(-20)
		{in: math.Pow(2, -20), want: 0x0010, acc: big.Exact},
		// 1 00000 0000000001 = -2^(-24)
		{in: -math.Pow(2, -24), want: 0x8001, acc: big.Exact},

		// 2^i
		{in: math.Pow(2, -25), want: 0x0000, acc: big.Below}, // 2^(-25)
		{in: math.Pow(2, -24), want: 0x0001, acc: big.Exact}, // 2^(-24)
		{in: math.Pow(2, -23), want: 0x0002, acc: big.Exact}, // 2^(-23)
		{in: math.Pow(2, -22), want: 0x0004, acc: big.Exact}, // 2^(-22)
		{in: math.Pow(2, -21), want: 0x0008, acc: big.Exact}, // 2^(-21)
		{in: math.Pow(2, -20), want: 0x0010, acc: big.Exact}, // 2^(-20)
		{in: math.Pow(2, -19), want: 0x0020, acc: big.Exact}, // 2^(-19)
		{in: math.Pow(2, -18), want: 0x0040, acc: big.Exact}, // 2^(-18)
		{in: math.Pow(2, -17), want: 0x0080, acc: big.Exact}, // 2^(-17)
		{in: math.Pow(2, -16), want: 0x0100, acc: big.Exact}, // 2^(-16)
		{in: math.Pow(2, -15), want: 0x0200, acc: big.Exact}, // 2^(-15)
		{in: math.Pow(2, -14), want: 0x0400, acc: big.Exact}, // 2^(-14)
		{in: math.Pow(2, -13), want: 0x0800, acc: big.Exact}, // 2^(-13)
		{in: math.Pow(2, -12), want: 0x0C00, acc: big.Exact}, // 2^(-12)
		{in: math.Pow(2, -11), want: 0x1000, acc: big.Exact}, // 2^(-11)
		{in: math.Pow(2, -10), want: 0x1400, acc: big.Exact}, // 2^(-10)
		{in: math.Pow(2, -9), want: 0x1800, acc: big.Exact},  // 2^(-9)
		{in: math.Pow(2, -8), want: 0x1C00, acc: big.Exact},  // 2^(-8)
		{in: math.Pow(2, -7), want: 0x2000, acc: big.Exact},  // 2^(-7)
		{in: math.Pow(2, -6), want: 0x2400, acc: big.Exact},  // 2^(-6)
		{in: math.Pow(2, -5), want: 0x2800, acc: big.Exact},  // 2^(-5)
		{in: math.Pow(2, -4), want: 0x2C00, acc: big.Exact},  // 2^(-4)
		{in: math.Pow(2, -3), want: 0x3000, acc: big.Exact},  // 2^(-3)
		{in: math.Pow(2, -2), want: 0x3400, acc: big.Exact},  // 2^(-2)
		{in: math.Pow(2, -1), want: 0x3800, acc: big.Exact},  // 2^(-1)
		{in: math.Pow(2, 0), want: 0x3C00, acc: big.Exact},   // 2^0
		{in: math.Pow(2, 1), want: 0x4000, acc: big.Exact},   // 2^1
		{in: math.Pow(2, 2), want: 0x4400, acc: big.Exact},   // 2^2
		{in: math.Pow(2, 3), want: 0x4800, acc: big.Exact},   // 2^3
		{in: math.Pow(2, 4), want: 0x4C00, acc: big.Exact},   // 2^4
		{in: math.Pow(2, 5), want: 0x5000, acc: big.Exact},   // 2^5
		{in: math.Pow(2, 6), want: 0x5400, acc: big.Exact},   // 2^6
		{in: math.Pow(2, 7), want: 0x5800, acc: big.Exact},   // 2^7
		{in: math.Pow(2, 8), want: 0x5C00, acc: big.Exact},   // 2^8
		{in: math.Pow(2, 9), want: 0x6000, acc: big.Exact},   // 2^9
		{in: math.Pow(2, 10), want: 0x6400, acc: big.Exact},  // 2^10
		{in: math.Pow(2, 11), want: 0x6800, acc: big.Exact},  // 2^11
		{in: math.Pow(2, 12), want: 0x6C00, acc: big.Exact},  // 2^12
		{in: math.Pow(2, 13), want: 0x7000, acc: big.Exact},  // 2^13
		{in: math.Pow(2, 14), want: 0x7400, acc: big.Exact},  // 2^14
		{in: math.Pow(2, 15), want: 0x7800, acc: big.Exact},  // 2^15
	}
	for _, g := range golden {
		f, acc := NewFromFloat64(g.in)
		got := f.Bits()
		x, _ := f.Float64()
		if g.want != got {
			t.Errorf("bits mismatch; expected 0x%04X (%v), got 0x%04X (%v)", g.want, g.in, got, x)
		}
		if g.acc != acc {
			t.Errorf("accuracy mismatch; expected %v (%v), got %v (%v)", g.acc, g.in, acc, x)
		}
	}
}
