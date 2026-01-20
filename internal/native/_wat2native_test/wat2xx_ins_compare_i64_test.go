// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestI64Eqz(t *testing.T) {
	const true, false = 1, 0
	if I64Eqz(0) == false {
		t.Error("0 should be Eqz")
	}
	if I64Eqz(1) == true {
		t.Error("1 should not be Eqz")
	}
	if I64Eqz(math.MinInt64) == true {
		t.Error("MinInt64 should not be Eqz")
	}
}
func TestI64Eq(t *testing.T) {
	const true, false = 1, 0
	if I64Eq(123456789012345, 123456789012345) == false {
		t.Error()
	}
	if I64Eq(1, 2) == true {
		t.Error()
	}
}
func TestI64Ne(t *testing.T) {
	const true, false = 1, 0
	if I64Ne(math.MaxInt64, math.MinInt64) == false {
		t.Error()
	}
	if I64Ne(100, 100) == true {
		t.Error()
	}
}
func TestI64Lt_s(t *testing.T) {
	const true, false = 1, 0
	// 有符号比较：负数小于正数
	if I64Lt_s(math.MinInt64, math.MaxInt64) == false {
		t.Error()
	}
	if I64Lt_s(-100, -50) == false {
		t.Error()
	}
	if I64Lt_s(0, -1) == true {
		t.Error("0 < -1 (signed) should be false")
	}
}
func TestI64Lt_u(t *testing.T) {
	const true, false = 1, 0
	// 无符号比较：0 最小，所有负数（高位为1）都极大
	if I64Lt_u(0, math.MaxInt64) == false {
		t.Error()
	}
	x := uint64(0xFFFFFFFFFFFFFFFF)
	if I64Lt_u(math.MaxInt64, int64(x)) == false {
		t.Error()
	}
	x = uint64(0xFFFFFFFFFFFFFFFF)
	if I64Lt_u(int64(x), 0) == true {
		t.Error()
	}
}
func TestI64Gt_s(t *testing.T) {
	const true, false = 1, 0
	if I64Gt_s(0, -1) == false {
		t.Error()
	}
	if I64Gt_s(math.MinInt64, math.MaxInt64) == true {
		t.Error()
	}
}
func TestI64Gt_u(t *testing.T) {
	const true, false = 1, 0
	// 无符号：0xFFFFFFFFFFFFFFFF > 0
	x := uint64(0xFFFFFFFFFFFFFFFF)
	if I64Gt_u(int64(x), 0) == false {
		t.Error()
	}
	x = uint64(0x8000000000000000)
	if I64Gt_u(math.MaxInt64, int64(x)) == true {
		t.Error()
	}
}
func TestI64Le_s(t *testing.T) {
	const true, false = 1, 0
	if I64Le_s(math.MinInt64, math.MinInt64) == false {
		t.Error()
	}
	if I64Le_s(-10, 0) == false {
		t.Error()
	}
}
func TestI64Le_u(t *testing.T) {
	const true, false = 1, 0
	x := uint64(0xFFFFFFFFFFFFFFFF)
	if I64Le_u(0, int64(x)) == false {
		t.Error()
	}
	if I64Le_u(int64(x), 0) == true {
		t.Error()
	}
}
func TestI64Ge_s(t *testing.T) {
	const true, false = 1, 0
	if I64Ge_s(0, -1) == false {
		t.Error()
	}
	if I64Ge_s(math.MaxInt64, math.MaxInt64) == false {
		t.Error()
	}
}
func TestI64Ge_u(t *testing.T) {
	const true, false = 1, 0
	// 无符号：最大值大于等于自己
	x := uint64(0xFFFFFFFFFFFFFFFF)
	maxU := int64(x)
	if I64Ge_u(maxU, maxU) == false {
		t.Error()
	}
	if I64Ge_u(maxU, 0) == false {
		t.Error()
	}
}
