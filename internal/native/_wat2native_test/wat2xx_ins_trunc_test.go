// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func assertPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	f()
}

func TestI32Trunc_f32_s(t *testing.T) {
	// 1. 正常转换
	if I32Trunc_f32_s(1.9) != 1 {
		t.Error()
	}
	if I32Trunc_f32_s(-1.9) != -1 {
		t.Error()
	}
	// 2. 边界检查 (MaxInt32 = 2147483647)
	if I32Trunc_f32_s(2147483520.0) != 2147483520 {
		t.Error()
	}
	// 3. 溢出与 NaN Trap
	assertPanic(t, func() { I32Trunc_f32_s(math.MaxFloat32) })
	assertPanic(t, func() { I32Trunc_f32_s(float32(math.NaN())) })
}
func TestI32Trunc_f32_u(t *testing.T) {
	if I32Trunc_f32_u(1.9) != 1 {
		t.Error()
	}
	// 4294967295.0 是 uint32_max
	if uint32(I32Trunc_f32_u(4294967040.0)) != 4294967040 {
		t.Error()
	}
	// 负数应报错
	assertPanic(t, func() { I32Trunc_f32_u(-1.0) })
}
func TestI32Trunc_f64_s(t *testing.T) {
	if I32Trunc_f64_s(2147483647.0) != 2147483647 {
		t.Error()
	}
	if I32Trunc_f64_s(-2147483648.0) != -2147483648 {
		t.Error()
	}
	assertPanic(t, func() { I32Trunc_f64_s(2147483648.0) })
}
func TestI32Trunc_f64_u(t *testing.T) {
	if uint32(I32Trunc_f64_u(4294967295.0)) != 4294967295 {
		t.Error()
	}
	assertPanic(t, func() { I32Trunc_f64_u(4294967296.0) })
}
func TestI64Trunc_f32_s(t *testing.T) {
	// 验证 64 位整数范围
	val := float32(9223372036854775807.0) // 可能会由于精度变为 2^63
	assertPanic(t, func() { I64Trunc_f32_s(val * 2) })
}
func TestI64Trunc_f32_u(t *testing.T) {
	if uint64(I64Trunc_f32_u(100.0)) != 100 {
		t.Error()
	}
	//assertPanic(t, func() { I64Trunc_f32_u(-0.1) })
}
func TestI64Trunc_f64_s(t *testing.T) {
	if I64Trunc_f64_s(1e18) != 1000000000000000000 {
		t.Error()
	}
	// 验证 NaN
	assertPanic(t, func() { I64Trunc_f64_s(math.NaN()) })
}
func TestI64Trunc_f64_u(t *testing.T) {
	// uint64 max 约为 1.84e19
	u64max_minus_eps := 18446744073709549568.0
	if uint64(I64Trunc_f64_u(u64max_minus_eps)) != 18446744073709549568 {
		t.Error()
	}
	assertPanic(t, func() { I64Trunc_f64_u(2e20) })
}
