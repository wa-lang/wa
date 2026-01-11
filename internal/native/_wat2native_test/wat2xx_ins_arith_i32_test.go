// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestI32Clz(t *testing.T) {
	tests := []struct {
		input    int32
		expected int32
	}{
		{0, 32},
		{1, 31},
		{-1, 0}, // 0xFFFFFFFF
		{0x000F0000, 12},
	}
	for _, tt := range tests {
		if got := I32Clz(tt.input); got != tt.expected {
			t.Errorf("I32Clz(%b) = %d; want %d", uint32(tt.input), got, tt.expected)
		}
	}
}

func TestI32Ctz(t *testing.T) {
	tests := []struct {
		input    uint32
		expected int32
	}{
		{0, 32},
		{1, 0},
		{0x80000000, 31},
		{0x000F0000, 16},
	}
	for _, tt := range tests {
		if got := I32Ctz(int32(tt.input)); got != tt.expected {
			t.Errorf("I32Ctz(%b) = %d; want %d", uint32(tt.input), got, tt.expected)
		}
	}
}
func TestI32Popcnt(t *testing.T) {
	if I32Popcnt(0) != 0 {
		t.Error("popcnt 0")
	}
	if I32Popcnt(-1) != 32 {
		t.Error("popcnt -1")
	}
	if I32Popcnt(0x10101010) != 4 {
		t.Error("popcnt hex")
	}
}

func TestI32Add_(t *testing.T) {
	if I32Add_(1, 2) != 3 {
		t.Error()
	}
	if I32Add_(math.MaxInt32, 1) != math.MinInt32 {
		t.Error("overflow")
	}
}

func TestI32Sub(t *testing.T) {
	if I32Sub(10, 3) != 7 {
		t.Error()
	}
	if I32Sub(math.MinInt32, 1) != math.MaxInt32 {
		t.Error("underflow")
	}
}

func TestI32Mul(t *testing.T) {
	if I32Mul(3, 4) != 12 {
		t.Error()
	}
	if I32Mul(-1, 5) != -5 {
		t.Error()
	}
	// 验证截断乘法
	// if I32Mul(0x12345678, 0x12345678) != 0x51E28100 {
	// 	t.Error("wrap")
	// }
}

func TestI32Div_s(t *testing.T) {
	if I32Div_s(10, 2) != 5 {
		t.Error()
	}
	if I32Div_s(-10, 2) != -5 {
		t.Error()
	}
	if I32Div_s(math.MinInt32, -1) != math.MinInt32 {
		t.Error("overflow case")
	}
}

func TestI32Div_u(t *testing.T) {
	// -1 as uint32 is 4294967295
	if I32Div_u(-1, 2) != 2147483647 {
		t.Error()
	}
}

func TestI32Rem_s(t *testing.T) {
	if I32Rem_s(10, 3) != 1 {
		t.Error()
	}
	if I32Rem_s(-10, 3) != -1 {
		t.Error()
	}
}

func TestI32Rem_u(t *testing.T) {
	// 4294967295 % 10 = 5
	if I32Rem_u(-1, 10) != 5 {
		t.Error()
	}
}

func TestI32And(t *testing.T) {
	if I32And(0xFF, 0x0F) != 0x0F {
		t.Error()
	}
}

func TestI32Or(t *testing.T) {
	if I32Or(0xF0, 0x0F) != 0xFF {
		t.Error()
	}
}

func TestI32Xor(t *testing.T) {
	if I32Xor(0xFF, 0xAA) != 0x55 {
		t.Error()
	}
}

func TestI32Shl(t *testing.T) {
	// Wasm: shift count is modulo 32
	if I32Shl(1, 1) != 2 {
		t.Error()
	}
	if I32Shl(1, 33) != 2 {
		t.Error("count % 32")
	}
}

func TestI32Shr_s(t *testing.T) {
	// Arithmetic shift: preserve sign
	if I32Shr_s(-2, 1) != -1 {
		t.Error()
	}
	if I32Shr_s(-1, 31) != -1 {
		t.Error()
	}
}

func TestI32Shr_u(t *testing.T) {
	// Logical shift: zero fill
	if I32Shr_u(-2, 1) != 0x7FFFFFFF {
		t.Error()
	}
}

func TestI32Rotl(t *testing.T) {
	// 0x80000000 rotate left 1 -> 1
	x := 0x80000000
	if I32Rotl(int32(x), 1) != 1 {
		t.Error()
	}
}

func TestI32Rotr(t *testing.T) {
	// 1 rotate right 1 -> 0x80000000
	x := 0x80000000
	if I32Rotr(1, 1) != int32(x) {
		t.Error()
	}
}
