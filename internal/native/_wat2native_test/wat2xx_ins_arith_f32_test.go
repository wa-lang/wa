// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestF32Abs(t *testing.T) {
	if F32Abs(-1.0) != 1.0 {
		t.Error()
	}
	if F32Abs(float32(math.Inf(-1))) != float32(math.Inf(1)) {
		t.Error()
	}
	// 验证 -0.0 变为 0.0
	res := F32Abs(float32(math.Copysign(0, -1)))
	if math.Signbit(float64(res)) {
		t.Error("Abs should clear sign bit of -0.0")
	}
}

func TestF32Neg(t *testing.T) {
	if F32Neg(1.0) != -1.0 {
		t.Error()
	}
	if F32Neg(-1.0) != 1.0 {
		t.Error()
	}
}
func TestF32Ceil(t *testing.T) {
	if F32Ceil(1.1) != 2.0 {
		t.Error()
	}
	if F32Ceil(-1.9) != -1.0 {
		t.Error()
	}
}
func TestF32Floor(t *testing.T) {
	if F32Floor(1.9) != 1.0 {
		t.Error()
	}
	if F32Floor(-1.1) != -2.0 {
		t.Error()
	}
}
func TestF32Trunc(t *testing.T) {
	if F32Trunc(1.9) != 1.0 {
		t.Error()
	}
	if F32Trunc(-1.9) != -1.0 {
		t.Error()
	}
}
func TestF32Nearest(t *testing.T) {
	if F32Nearest(0.5) != 0.0 {
		t.Errorf("0.5 -> 0.0, got %f", F32Nearest(0.5))
	}
	if F32Nearest(1.5) != 2.0 {
		t.Errorf("1.5 -> 2.0, got %f", F32Nearest(1.5))
	}
	if F32Nearest(2.5) != 2.0 {
		t.Errorf("2.5 -> 2.0, got %f", F32Nearest(2.5))
	}
}
func TestF32Sqrt(t *testing.T) {
	if F32Sqrt(4.0) != 2.0 {
		t.Error()
	}
	if !math.IsNaN(float64(F32Sqrt(-1.0))) {
		t.Error("Sqrt of negative should be NaN")
	}
}
func TestF32Add(t *testing.T) {
	if F32Add(1.5, 2.5) != 4.0 {
		t.Error()
	}
}
func TestF32Sub(t *testing.T) {
	if F32Sub(5.0, 2.5) != 2.5 {
		t.Error()
	}
}
func TestF32Mul(t *testing.T) {
	if F32Mul(2.0, 3.5) != 7.0 {
		t.Error()
	}
}
func TestF32Div(t *testing.T) {
	if F32Div(10.0, 2.0) != 5.0 {
		t.Error()
	}
}
func TestF32Min(t *testing.T) {
	pos0 := float32(0.0)
	neg0 := float32(math.Copysign(0, -1))

	// 1. 规范要求 min(0, -0) = -0
	res := F32Min(pos0, neg0)
	if !math.Signbit(float64(res)) {
		t.Error("Min(0, -0) should be -0")
	}

	// 2. 规范要求 NaN 传播
	if !math.IsNaN(float64(F32Min(pos0, float32(math.NaN())))) {
		t.Error("Min with NaN should be NaN")
	}

	if F32Min(1.0, 2.0) != 1.0 {
		t.Error()
	}
}
func TestF32Max(t *testing.T) {
	pos0 := float32(0.0)
	neg0 := float32(math.Copysign(0, -1))

	// 1. 规范要求 max(0, -0) = 0
	res := F32Max(pos0, neg0)
	if math.Signbit(float64(res)) {
		t.Error("Max(0, -0) should be 0")
	}

	if F32Max(-10.0, 5.0) != 5.0 {
		t.Error()
	}
}
func TestF32Copysign(t *testing.T) {
	if F32Copysign(1.0, -2.0) != -1.0 {
		t.Error()
	}
	if F32Copysign(-5.0, 0.0) != 5.0 {
		t.Error()
	}
}
