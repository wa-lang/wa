// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

import "math"

func I32Reinterpret_f32(v float32) int32 {
	return int32(math.Float32bits(v))
}
func I64Reinterpret_f64(v float64) int64 {
	return int64(math.Float64bits(v))
}
func F32Reinterpret_i32(v int32) float32 {
	return math.Float32frombits(uint32(v))
}
func F64Reinterpret_i64(v int64) float64 {
	return math.Float64frombits(uint64(v))
}
