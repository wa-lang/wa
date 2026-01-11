// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestI32Reinterpret_f32(t *testing.T) {
	// 构建一个特殊的位模式：0xBF800000 (这是 float32 的 -1.0)
	bits := uint32(0xBF800000)
	f := math.Float32frombits(bits)

	got := I32Reinterpret_f32(f)
	if uint32(got) != bits {
		t.Errorf("I32Reinterpret_f32 failed: expected 0x%X, got 0x%X", bits, got)
	}

	// 测试 NaN 的位模式保留
	nanBits := uint32(0x7FC00000)
	nanF := math.Float32frombits(nanBits)
	if uint32(I32Reinterpret_f32(nanF)) != nanBits {
		t.Error("NaN payload was not preserved during reinterpret")
	}
}
func TestI64Reinterpret_f64(t *testing.T) {
	// 0x400921FB54442D18 是 float64 的 Pi (3.141592653589793)
	bits := uint64(0x400921FB54442D18)
	f := math.Float64frombits(bits)

	got := I64Reinterpret_f64(f)
	if uint64(got) != bits {
		t.Errorf("I64Reinterpret_f64 failed: expected 0x%X, got 0x%X", bits, got)
	}
}
func TestF32Reinterpret_i32(t *testing.T) {
	// 0x3F800000 解释为 float32 应为 1.0
	bits := int32(0x3F800000)
	got := F32Reinterpret_i32(bits)

	if got != 1.0 {
		t.Errorf("F32Reinterpret_i32 failed: expected 1.0, got %f", got)
	}

	// 验证符号位：0x80000000 应为 -0.0
	x := uint32(0x80000000)
	negZeroBits := int32(x)
	res := F32Reinterpret_i32(negZeroBits)
	if res != 0.0 || !math.Signbit(float64(res)) {
		t.Error("Should be negative zero")
	}
}
func TestF64Reinterpret_i64(t *testing.T) {
	x := uint64(0xBFF0000000000000)
	// 0xBFF0000000000000 解释为 float64 应为 -1.0
	bits := int64(x)
	got := F64Reinterpret_i64(bits)

	if got != -1.0 {
		t.Errorf("F64Reinterpret_i64 failed: expected -1.0, got %f", got)
	}
}
