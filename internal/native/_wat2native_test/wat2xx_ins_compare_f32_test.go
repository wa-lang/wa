// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

package main

import (
	"math"
	"testing"
)

func TestF32Eq(t *testing.T) {
	const true, false = 1, 0
	var neg0 = math.Float32frombits(0x80000000)

	nan := float32(math.NaN())
	if F32Eq(1.0, 1.0) == false {
		t.Error("1.0 == 1.0")
	}
	if F32Eq(1.0, 2.0) == true {
		t.Error("1.0 != 2.0")
	}
	// 规范：NaN 参与比较永远为 false
	if F32Eq(nan, nan) == true {
		t.Error("NaN == NaN should be false")
	}
	// 规范：0.0 == -0.0
	if F32Eq(0.0, neg0) == false {
		t.Error("0.0 == -0.0 should be true")
	}
}
func TestF32Ne(t *testing.T) {
	const true, false = 1, 0

	nan := float32(math.NaN())
	if F32Ne(1.0, 2.0) == false {
		t.Error("1.0 != 2.0")
	}
	if F32Ne(1.0, 1.0) == true {
		t.Error("1.0 == 1.0")
	}
	// 规范：NaN != 任何值 结果为 true
	if F32Ne(nan, 1.0) == false {
		t.Error("NaN != 1.0 should be true")
	}
	if F32Ne(nan, nan) == false {
		t.Error("NaN != NaN should be true")
	}
}
func TestF32Lt(t *testing.T) {
	const true, false = 1, 0
	var neg0 = math.Float32frombits(0x80000000)

	if F32Lt(1.0, 2.0) == false {
		t.Error()
	}
	if F32Lt(2.0, 1.0) == true {
		t.Error()
	}
	// NaN 比较
	if F32Lt(float32(math.NaN()), 1.0) == true {
		t.Error("NaN < 1.0 should be false")
	}
	// 0.0 与 -0.0
	if F32Lt(neg0, 0.0) == true {
		t.Error("-0.0 < 0.0 should be false")
	}
}
func TestF32Gt(t *testing.T) {
	const true, false = 1, 0
	if F32Gt(2.0, 1.0) == false {
		t.Error()
	}
	if F32Gt(float32(math.NaN()), 1.0) == true {
		t.Error()
	}
}
func TestF32Le(t *testing.T) {
	const true, false = 1, 0
	var neg0 = math.Float32frombits(0x80000000)

	if F32Le(1.0, 1.0) == false {
		t.Error()
	}
	if F32Le(1.0, 2.0) == false {
		t.Error()
	}
	// 0.0 <= -0.0
	if F32Le(0.0, neg0) == false {
		t.Error("0.0 <= -0.0 should be true")
	}
	// NaN
	if F32Le(float32(math.NaN()), 1.0) == true {
		t.Error()
	}
}
func TestF32Ge(t *testing.T) {
	const true, false = 1, 0
	if F32Ge(1.0, 1.0) == false {
		t.Error()
	}
	if F32Ge(2.0, 1.0) == false {
		t.Error()
	}
	// NaN
	if F32Ge(float32(math.NaN()), 1.0) == true {
		t.Error()
	}
}
