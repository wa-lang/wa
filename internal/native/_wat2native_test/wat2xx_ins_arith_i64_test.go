// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestI64Clz(t *testing.T) {
	tests := []struct {
		input    int64
		expected int64
	}{
		{0, 64},
		{1, 63},
		{-1, 0}, // 所有位均为 1
		{0x00000000FFFFFFFF, 32},
	}
	for _, tt := range tests {
		if got := I64Clz(tt.input); got != tt.expected {
			t.Errorf("I64Clz(%x) = %d; want %d", uint64(tt.input), got, tt.expected)
		}
	}
}
func TestI64Ctz(t *testing.T) {
	tests := []struct {
		input    int64
		expected int64
	}{
		{0, 64},
		{1, 0},
		{math.MinInt64, 63}, // 0x8000...0000
		{0x0000FFFF00000000, 32},
	}
	for _, tt := range tests {
		if got := I64Ctz(tt.input); got != tt.expected {
			t.Errorf("I64Ctz(%x) = %d; want %d", uint64(tt.input), got, tt.expected)
		}
	}
}
func TestI64Popcnt(t *testing.T) {
	if I64Popcnt(0) != 0 {
		t.Error()
	}
	if I64Popcnt(-1) != 64 {
		t.Error()
	}
	if I64Popcnt(0x1234567812345678) != 26 {
		t.Error()
	}
}
func TestI64Add(t *testing.T) {
	if I64Add(1, 1) != 2 {
		t.Error()
	}
	if I64Add(math.MaxInt64, 1) != math.MinInt64 {
		t.Error("I64 overflow")
	}
}
func TestI64Sub(t *testing.T) {
	if I64Sub(0, 1) != -1 {
		t.Error()
	}
	if I64Sub(math.MinInt64, 1) != math.MaxInt64 {
		t.Error("I64 underflow")
	}
}
func TestI64Mul(t *testing.T) {
	if I64Mul(100, 200) != 20000 {
		t.Error()
	}
	// 验证 64 位截断
	res := I64Mul(0x7FFFFFFFFFFFFFFF, 2)
	if res != -2 {
		t.Error("I64 mul wrap")
	}
}
func TestI64Div_s(t *testing.T) {
	if I64Div_s(-100, 10) != -10 {
		t.Error()
	}
	if I64Div_s(math.MinInt64, -1) != math.MinInt64 {
		t.Error("I64 div overflow")
	}
}
func TestI64Div_u(t *testing.T) {
	// -1 as uint64 is 18446744073709551615
	if I64Div_u(-1, 2) != 9223372036854775807 {
		t.Error()
	}
}
func TestI64Rem_s(t *testing.T) {
	if I64Rem_s(-10, 3) != -1 {
		t.Error()
	}
	if I64Rem_s(math.MinInt64, -1) != 0 {
		t.Error()
	}
}
func TestI64Rem_u(t *testing.T) {
	if I64Rem_u(-1, 10) != 5 {
		t.Error()
	}
}
func TestI64And(t *testing.T) {
	if I64And(0xFFFF, 0x00FF) != 0x00FF {
		t.Error()
	}
}
func TestI64Or(t *testing.T) {
	if I64Or(0x0F, 0xF0) != 0xFF {
		t.Error()
	}
}
func TestI64Xor(t *testing.T) {
	if I64Xor(0xAAAA, 0x5555) != 0xFFFF {
		t.Error()
	}
}
func TestI64Shl(t *testing.T) {
	if I64Shl(1, 1) != 2 {
		t.Error()
	}
	if I64Shl(1, 65) != 2 {
		t.Error("count % 64")
	}
}
func TestI64Shr_s(t *testing.T) {
	if I64Shr_s(-1, 63) != -1 {
		t.Error()
	}
	if I64Shr_s(math.MinInt64, 1) != -4611686018427387904 {
		t.Error()
	}
}
func TestI64Shr_u(t *testing.T) {
	if I64Shr_u(-1, 63) != 1 {
		t.Error()
	}
}
func TestI64Rotl(t *testing.T) {
	val := uint64(1) << 63
	if I64Rotl(int64(val), 1) != 1 {
		t.Errorf("Rotl fail: got %d", I64Rotl(int64(val), 1))
	}
}
func TestI64Rotr(t *testing.T) {
	val := uint64(1) << 63
	expected := int64(val)
	if I64Rotr(1, 1) != expected {
		t.Error()
	}
}
