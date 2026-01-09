// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestGlobalI32(t *testing.T) {
	tests := []int32{
		0,
		1,
		-1,
		math.MaxInt32,
		math.MinInt32,
		0x12345678,
	}
	for _, v := range tests {
		GlobalSetI32(v)
		if got := GlobalGetI32(); got != v {
			t.Errorf("Global I32 读写不匹配: set %d, got %d", v, got)
		}
	}
}

func TestGlobalI64(t *testing.T) {
	tests := []int64{
		0,
		1,
		-1,
		math.MaxInt64,
		math.MinInt64,
		0x1234567890ABCDEF,
	}
	for _, v := range tests {
		GlobalSetI64(v)
		if got := GlobalGetI64(); got != v {
			t.Errorf("Global I64 读写不匹配: set %d, got %d", v, got)
		}
	}
}

func TestGlobalF32(t *testing.T) {
	tests := []float32{
		0,
		3.1415927,
		-1.0,
		float32(math.MaxFloat32),
		float32(math.SmallestNonzeroFloat32),
	}
	for _, v := range tests {
		GlobalSetF32(v)
		got := GlobalGetF32()

		// 浮点数比较通常需要注意精度
		// 但在这种直接搬运的场景下, 二进制位应该完全一致

		if got != v {
			t.Errorf("Global F32 读写不匹配: set %v, got %v", v, got)
		}
	}
}

func TestGlobalF64(t *testing.T) {
	tests := []float64{
		0,
		math.Pi,
		-1.0,
		math.MaxFloat64,
		math.SmallestNonzeroFloat64,
	}
	for _, v := range tests {
		GlobalSetF64(v)
		got := GlobalGetF64()
		if got != v {
			t.Errorf("Global F64 读写不匹配: set %v, got %v", v, got)
		}
	}
}

func TestGlobalInterleaved(t *testing.T) {
	GlobalSetI32(42)
	GlobalSetF64(math.E)
	GlobalSetI64(0xDEADBEEF)

	if GlobalGetI32() != 42 {
		t.Error("I32 value corrupted after other sets")
	}
	if GlobalGetF64() != math.E {
		t.Error("F64 value corrupted after other sets")
	}
}
