// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestI32Eqz(t *testing.T) {
	const true, false = 1, 0
	if I32Eqz(0) == false {
		t.Error("0 should be Eqz")
	}
	if I32Eqz(1) == true {
		t.Error("1 should not be Eqz")
	}
	if I32Eqz(-1) == true {
		t.Error("-1 should not be Eqz")
	}
}
func TestI32Eq(t *testing.T) {
	const true, false = 1, 0
	if I32Eq(10, 10) == false {
		t.Error()
	}
	if I32Eq(10, 20) == true {
		t.Error()
	}
}
func TestI32Ne(t *testing.T) {
	const true, false = 1, 0
	if I32Ne(10, 20) == false {
		t.Error()
	}
	if I32Ne(10, 10) == true {
		t.Error()
	}
}
func TestI32Lt_s(t *testing.T) {
	const true, false = 1, 0
	// 有符号比较：-1 < 0
	if I32Lt_s(-1, 0) == false {
		t.Error("-1 should be < 0 (signed)")
	}
	if I32Lt_s(1, -1) == true {
		t.Error("1 should not be < -1 (signed)")
	}
}
func TestI32Lt_u(t *testing.T) {
	const true, false = 1, 0
	// 无符号比较：-1 (0xFFFFFFFF) > 0
	if I32Lt_u(-1, 0) == true {
		t.Error("-1 (uint32_max) should not be < 0 (unsigned)")
	}
	if I32Lt_u(0, -1) == false {
		t.Error("0 should be < -1 (unsigned)")
	}
}
func TestI32Gt_s(t *testing.T) {
	const true, false = 1, 0
	if I32Gt_s(0, -1) == false {
		t.Error("0 should be > -1 (signed)")
	}
	if I32Gt_s(math.MinInt32, 0) == true {
		t.Error("min_int should not be > 0")
	}
}
func TestI32Gt_u(t *testing.T) {
	const true, false = 1, 0
	// 无符号：0x80000000 > 0x7FFFFFFF
	if I32Gt_u(math.MinInt32, math.MaxInt32) == false {
		t.Error("unsigned: 0x80000000 should be > 0x7FFFFFFF")
	}
}
func TestI32Le_s(t *testing.T) {
	const true, false = 1, 0
	if I32Le_s(-1, -1) == false {
		t.Error()
	}
	if I32Le_s(-2, -1) == false {
		t.Error()
	}
	if I32Le_s(0, -1) == true {
		t.Error()
	}
}
func TestI32Le_u(t *testing.T) {
	const true, false = 1, 0
	if I32Le_u(0, 0) == false {
		t.Error()
	}
	if I32Le_u(10, -1) == false {
		t.Error("10 should be <= 0xFFFFFFFF")
	}
}
func TestI32Ge_s(t *testing.T) {
	const true, false = 1, 0
	if I32Ge_s(math.MaxInt32, math.MinInt32) == false {
		t.Error()
	}
}
func TestI32Ge_u(t *testing.T) {
	const true, false = 1, 0
	// 无符号：最小负数（最高位1）大于最大正数（最高位0）
	if I32Ge_u(math.MinInt32, math.MaxInt32) == false {
		t.Error()
	}
}
