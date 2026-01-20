// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

// 定义位模式常量，方便构建特殊值
const (
	bitsPos0 uint64 = 0x0000000000000000
	bitsNeg0 uint64 = 0x8000000000000000
	bitsNaN  uint64 = 0x7FF8000000000001 // 一个标准的 Quiet NaN
)

func TestF64Eq(t *testing.T) {
	const true, false = 1, 0
	p0 := math.Float64frombits(bitsPos0)
	n0 := math.Float64frombits(bitsNeg0)
	nan := math.Float64frombits(bitsNaN)

	// 1. 验证 0.0 == -0.0
	if F64Eq(p0, n0) == false {
		t.Error("F64Eq(0.0, -0.0) should be true per Wasm spec")
	}
	// 2. 验证 NaN != NaN
	if F64Eq(nan, nan) == true {
		t.Error("F64Eq(NaN, NaN) should be false")
	}
	// 3. 普通数值
	if F64Eq(1.23, 1.23) == false {
		t.Error("F64Eq(1.23, 1.23) failed")
	}
}
func TestF64Ne(t *testing.T) {
	const true, false = 1, 0
	p0 := math.Float64frombits(bitsPos0)
	n0 := math.Float64frombits(bitsNeg0)
	nan := math.Float64frombits(bitsNaN)

	// 1. 验证 0.0 != -0.0 为假
	if F64Ne(p0, n0) == true {
		t.Error("F64Ne(0.0, -0.0) should be false")
	}
	// 2. 验证 NaN != 任何值（包括自己）为真
	if F64Ne(nan, nan) == false {
		t.Error("F64Ne(NaN, NaN) should be true")
	}
	if F64Ne(nan, 1.0) == false {
		t.Error("F64Ne(NaN, 1.0) should be true")
	}
}
func TestF64Lt(t *testing.T) {
	const true, false = 1, 0
	p0 := math.Float64frombits(bitsPos0)
	n0 := math.Float64frombits(bitsNeg0)
	nan := math.Float64frombits(bitsNaN)

	// 1. 验证 -0.0 不小于 0.0
	if F64Lt(n0, p0) == true {
		t.Error("F64Lt(-0.0, 0.0) should be false")
	}
	// 2. 验证 NaN 比较结果
	if F64Lt(nan, 100.0) == true {
		t.Error("F64Lt(NaN, 100.0) should be false")
	}
	// 3. 正常逻辑
	if F64Lt(-1.0, 0.0) == false {
		t.Error("F64Lt(-1.0, 0.0) failed")
	}
}
func TestF64Gt(t *testing.T) {
	const true, false = 1, 0
	p0 := math.Float64frombits(bitsPos0)
	n0 := math.Float64frombits(bitsNeg0)
	nan := math.Float64frombits(bitsNaN)

	if F64Gt(p0, n0) == true {
		t.Error("F64Gt(0.0, -0.0) should be false")
	}
	if F64Gt(nan, -math.MaxFloat64) == true {
		t.Error("F64Gt(NaN, -Inf) should be false")
	}
}
func TestF64Le(t *testing.T) {
	const true, false = 1, 0
	p0 := math.Float64frombits(bitsPos0)
	n0 := math.Float64frombits(bitsNeg0)
	nan := math.Float64frombits(bitsNaN)

	// 1. 验证 0.0 <= -0.0 为真
	if F64Le(p0, n0) == false {
		t.Error("F64Le(0.0, -0.0) should be true")
	}
	// 2. 验证无穷大
	inf := math.Inf(1)
	if F64Le(inf, inf) == false {
		t.Error("F64Le(+Inf, +Inf) should be true")
	}
	// 3. NaN 仍然为假
	if F64Le(nan, inf) == true {
		t.Error("F64Le(NaN, Inf) should be false")
	}
}
func TestF64Ge(t *testing.T) {
	const true, false = 1, 0
	p0 := math.Float64frombits(bitsPos0)
	n0 := math.Float64frombits(bitsNeg0)
	nan := math.Float64frombits(bitsNaN)

	if F64Ge(n0, p0) == false {
		t.Error("F64Ge(-0.0, 0.0) should be true")
	}
	if F64Ge(nan, nan) == true {
		t.Error("F64Ge(NaN, NaN) should be false")
	}
}
