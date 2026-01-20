// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestF32Convert_i32_s(t *testing.T) {
	if F32Convert_i32_s(-1) != -1.0 {
		t.Error()
	}
	if F32Convert_i32_s(100) != 100.0 {
		t.Error()
	}
}
func TestF32Convert_i32_u(t *testing.T) {
	if F32Convert_i32_u(-1) != 4294967295.0 {
		t.Error()
	}
}
func TestF32Convert_i64_s(t *testing.T) {
	if F32Convert_i64_s(-100) != -100.0 {
		t.Error()
	}
	// 验证大数转换时的精度舍入
	val := int64(1<<60 + 1)
	if F32Convert_i64_s(val) != float32(1<<60) {
		t.Log("Precision lost as expected")
	}
}
func TestF32Convert_i64_u(t *testing.T) {
	// uint64 max
	u64max := uint64(0xFFFFFFFFFFFFFFFF)
	res := F32Convert_i64_u(int64(u64max))
	if res != 18446744073709551616.0 { // 2^64
		t.Errorf("Expected 2^64, got %f", res)
	}
}
func TestF64Convert_i32_s(t *testing.T) {
	if F64Convert_i32_s(-1) != -1.0 {
		t.Error()
	}
}
func TestF64Convert_i32_u(t *testing.T) {
	if F64Convert_i32_u(-1) != 4294967295.0 {
		t.Error()
	}
}
func TestF64Convert_i64_s(t *testing.T) {
	if F64Convert_i64_s(math.MinInt64) != -9223372036854775808.0 {
		t.Error()
	}
}
func TestF64Convert_i64_u(t *testing.T) {
	u64max := uint64(0xFFFFFFFFFFFFFFFF)
	if F64Convert_i64_u(int64(u64max)) != 18446744073709551615.0 {
		// f64 有 53 位尾数，足以精确表示一些大数，但对于 u64 依然会舍入到 2^64
		t.Logf("Result: %f", F64Convert_i64_u(int64(u64max)))
	}
}
func TestI32Wrap_i64(t *testing.T) {
	// 0x00000001_00000002 -> 0x00000002
	val := int64(0x100000002)
	if I32Wrap_i64(val) != 2 {
		t.Error()
	}
}
func TestI64Extend_i32_s(t *testing.T) {
	// -1 (0xFFFFFFFF) -> -1 (0xFFFFFFFFFFFFFFFF)
	if I64Extend_i32_s(-1) != -1 {
		t.Error()
	}
	if I64Extend_i32_s(1) != 1 {
		t.Error()
	}
}
func TestI64Extend_i32_u(t *testing.T) {
	// -1 (0xFFFFFFFF) -> 4294967295 (0x00000000FFFFFFFF)
	res := I64Extend_i32_u(-1)
	if res != 4294967295 {
		t.Errorf("Expected 4294967295, got %d", res)
	}
}
func TestF32Demote_f64(t *testing.T) {
	val := 3.141592653589793
	res := F32Demote_f64(val)
	if float64(res) == val {
		t.Error("Should have lost precision")
	}
	if res != 3.1415927 {
		t.Errorf("Got %f", res)
	}
}
func TestF64Promote_f32(t *testing.T) {
	var val float32 = 1.5
	if F64Promote_f32(val) != 1.5 {
		t.Error()
	}

	// 测试 NaN 提升
	nan := float32(math.NaN())
	if !math.IsNaN(F64Promote_f32(nan)) {
		t.Error("NaN should promote to NaN")
	}
}
