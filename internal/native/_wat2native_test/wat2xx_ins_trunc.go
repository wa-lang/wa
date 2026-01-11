// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

import (
	"math"
)

func I32Trunc_f32_s(v float32) int32 {
	if math.IsNaN(float64(v)) {
		panic("invalid conversion to integer")
	}
	if v >= math.MinInt32-1 && v < math.MaxInt32+1 {
		return int32(v)
	}
	panic("integer overflow")
}
func I32Trunc_f32_u(v float32) int32 {
	if math.IsNaN(float64(v)) {
		panic("invalid conversion to integer")
	}
	if v > -1.0 && v < math.MaxUint32+1 {
		return int32(uint32(v))
	}
	panic("integer overflow")
}
func I32Trunc_f64_s(v float64) int32 {
	if math.IsNaN(v) {
		panic("invalid conversion to integer")
	}
	if v > math.MinInt32-1 && v < math.MaxInt32+1 {
		return int32(v)
	}
	panic("integer overflow")
}
func I32Trunc_f64_u(v float64) int32 {
	if math.IsNaN(v) {
		panic("invalid conversion to integer")
	}
	if v > -1.0 && v < math.MaxUint32+1 {
		return int32(uint32(v))
	}
	panic("integer overflow")
}

func I64Trunc_f32_s(v float32) int64 {
	if math.IsNaN(float64(v)) {
		panic("invalid conversion to integer")
	}
	if v >= float32(math.MinInt64) && v < -float32(math.MinInt64) {
		return int64(v)
	}
	panic("integer overflow")
}
func I64Trunc_f32_u(v float32) int64 {
	if math.IsNaN(float64(v)) {
		panic("invalid conversion to integer")
	}
	if v > -1.0 && v < -2.0*float32(math.MinInt64) {
		return int64(uint64(v))
	}
	panic("integer overflow")
}
func I64Trunc_f64_s(v float64) int64 {
	if math.IsNaN(v) {
		panic("invalid conversion to integer")
	}
	if v >= math.MinInt64 && v < -math.MinInt64 {
		return int64(v)
	}
	panic("integer overflow")
}
func I64Trunc_f64_u(v float64) int64 {
	if math.IsNaN(v) {
		panic("invalid conversion to integer")
	}
	if v > -1.0 && v < 18446744073709551616.0 { // 2^64
		return int64(uint64(v))
	}
	panic("integer overflow")
}
