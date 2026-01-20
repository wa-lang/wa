// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestF64Abs(t *testing.T) {
	if F64Abs(-3.14) != 3.14 {
		t.Error()
	}
	if F64Abs(math.Inf(-1)) != math.Inf(1) {
		t.Error()
	}
	// 验证 -0.0 的符号位被清除
	res := F64Abs(math.Copysign(0, -1))
	if math.Signbit(res) {
		t.Error("Abs failed to clear sign bit of -0.0")
	}
}
func TestF64Neg(t *testing.T) {
	if F64Neg(3.14) != -3.14 {
		t.Error()
	}
	if F64Neg(-3.14) != 3.14 {
		t.Error()
	}
}
func TestF64Ceil(t *testing.T) {
	if F64Ceil(1.1) != 2.0 {
		t.Error()
	}
	if F64Ceil(-1.9) != -1.0 {
		t.Error()
	}
}
func TestF64Floor(t *testing.T) {
	if F64Floor(1.9) != 1.0 {
		t.Error()
	}
	if F64Floor(-1.1) != -2.0 {
		t.Error()
	}
}
func TestF64Trunc(t *testing.T) {
	if F64Trunc(1.9) != 1.0 {
		t.Error()
	}
	if F64Trunc(-1.9) != -1.0 {
		t.Error()
	}
}
func TestF64Nearest(t *testing.T) {
	if F64Nearest(0.5) != 0.0 {
		t.Errorf("0.5 -> 0.0, got %f", F64Nearest(0.5))
	}
	if F64Nearest(1.5) != 2.0 {
		t.Errorf("1.5 -> 2.0, got %f", F64Nearest(1.5))
	}
	if F64Nearest(2.5) != 2.0 {
		t.Errorf("2.5 -> 2.0, got %f", F64Nearest(2.5))
	}
	if F64Nearest(3.5) != 4.0 {
		t.Errorf("3.5 -> 4.0, got %f", F64Nearest(3.5))
	}
}
func TestF64Sqrt(t *testing.T) {
	if F64Sqrt(9.0) != 3.0 {
		t.Error()
	}
	if !math.IsNaN(F64Sqrt(-1.0)) {
		t.Error("Sqrt(-1) should be NaN")
	}
}
func TestF64Add(t *testing.T) {
	if F64Add(1e15, 1.0) != 1000000000000001.0 {
		t.Error()
	}
}
func TestF64Sub(t *testing.T) {
	a, b := 10.0, 3.14
	c := a - b
	if F64Sub(a, b) != c {
		t.Error()
	}
}
func TestF64Mul(t *testing.T) {
	if F64Mul(2.5, 4.0) != 10.0 {
		t.Error()
	}
}
func TestF64Div(t *testing.T) {
	if F64Div(1.0, 2.0) != 0.5 {
		t.Error()
	}
	if !math.IsInf(F64Div(1.0, 0.0), 1) {
		t.Error("Divide by zero should be Inf")
	}
}
func TestF64Min(t *testing.T) {
	pos0 := 0.0
	neg0 := math.Copysign(0, -1)

	// Wasm 规范: min(0.0, -0.0) == -0.0
	res := F64Min(pos0, neg0)
	if !math.Signbit(res) {
		t.Error("Min(0, -0) should be -0")
	}

	// NaN 传播验证
	if !math.IsNaN(F64Min(math.NaN(), 1.0)) {
		t.Error("Min with NaN should return NaN")
	}
}
func TestF64Max(t *testing.T) {
	pos0 := 0.0
	neg0 := math.Copysign(0, -1)

	// Wasm 规范: max(0.0, -0.0) == 0.0
	res := F64Max(pos0, neg0)
	if math.Signbit(res) {
		t.Error("Max(0, -0) should be 0")
	}

	if F64Max(-100.0, -200.0) != -100.0 {
		t.Error()
	}
}
func TestF64Copysign(t *testing.T) {
	res := F64Copysign(1.23, math.Copysign(0, -1))
	if res != -1.23 || !math.Signbit(res) {
		t.Errorf("Copysign failed, got %f", res)
	}
}
