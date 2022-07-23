package bfloat

import (
	"math"
	"testing"
)

func TestNewFromBits(t *testing.T) {
	golden := []struct {
		bits uint16
		want float64
	}{
		// Special numbers.
		// 0 00000000 0000000 = 0
		{bits: 0, want: 0},
		// 1 00000000 0000000 = -0
		{bits: 0x8000, want: 1. / math.Inf(-1)},
		// 0 11111111 0000000 = +Inf
		{bits: 0x7f80, want: math.Inf(1)},
		// 1 11111111 0000000 = -Inf
		{bits: 0xff80, want: math.Inf(-1)},

		// 0 11111111 0000001 = +NaN
		{bits: 0x7f81, want: math.NaN()},
		// 1 11111111 0000001 = -NaN
		{bits: 0xff81, want: -math.NaN()},

		// from: https://en.wikipedia.org/wiki/Bfloat16_floating-point_format#Examples
		{bits: 0x3f80, want: 1},
		{bits: 0xc000, want: -2},
		{bits: 0x4049, want: 3.140625},
		{bits: 0x3eab, want: 0.333984375},
	}
	for _, g := range golden {
		f := NewFromBits(g.bits)
		b, isNan := f.Big()
		got, _ := b.Float64()
		if isNan {
			got = g.want
		}
		wantBits := math.Float64bits(g.want)
		gotBits := math.Float64bits(got)
		if wantBits != gotBits {
			t.Errorf("0x%04X: number mismatch; expected 0x%016X (%v), got 0x%016X (%v)", g.bits, wantBits, g.want, gotBits, got)
		}
	}
}
