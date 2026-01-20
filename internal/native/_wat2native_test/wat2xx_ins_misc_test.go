// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestSelectI32(t *testing.T) {
	// 1. 条件为 1 (非0)，选择 a
	if SelectI32(1, 10, 20) != 10 {
		t.Error("SelectI32(1, 10, 20) should return 10")
	}
	// 2. 条件为 0，选择 b
	if SelectI32(0, 10, 20) != 20 {
		t.Error("SelectI32(0, 10, 20) should return 20")
	}
	// 3. 边界情况：非 0 的负数条件，选择 a
	if SelectI32(-1, 10, 20) != 10 {
		t.Error("SelectI32(-1, 10, 20) should return 10")
	}
}
func TestSelectI64(t *testing.T) {
	var tb = uint64(0x99AABBCCDDEEFF00)
	var a int64 = 0x1122334455667788
	var b int64 = int64(tb)

	if SelectI64(1, a, b) != a {
		t.Error()
	}
	if SelectI64(0, a, b) != b {
		t.Error()
	}
	// 验证即使是很大的负数，只要不为 0 就选 a
	if SelectI64(math.MinInt32, a, b) != a {
		t.Error()
	}
}
func TestSelectF32(t *testing.T) {
	var valA float32 = 3.14
	var valB float32 = 2.71

	if SelectF32(1, valA, valB) != valA {
		t.Error()
	}
	if SelectF32(0, valA, valB) != valB {
		t.Error()
	}

	// 测试 NaN 作为操作数被选择的情况
	nan := float32(math.NaN())
	res := SelectF32(1, nan, valB)
	if !math.IsNaN(float64(res)) {
		t.Error("Should select NaN")
	}
}
func TestSelectF64(t *testing.T) {
	var valA float64 = 1.23456789
	var valB float64 = 9.87654321

	if SelectF64(1, valA, valB) != valA {
		t.Error()
	}
	if SelectF64(0, valA, valB) != valB {
		t.Error()
	}
}
