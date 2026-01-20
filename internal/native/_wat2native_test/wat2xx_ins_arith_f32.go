// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

import "math"

func F32Abs(v float32) float32 {
	return math.Float32frombits(math.Float32bits(v) &^ (1 << 31))
}
func F32Neg(v float32) float32 {
	return math.Float32frombits(math.Float32bits(v) ^ (1 << 31))
}
func F32Ceil(v float32) float32 {
	return float32(math.Ceil(float64(v)))
}
func F32Floor(v float32) float32 {
	return float32(math.Floor(float64(v)))
}
func F32Trunc(v float32) float32 {
	return float32(math.Trunc(float64(v)))
}
func F32Nearest(v float32) float32 {
	return float32(math.RoundToEven(float64(v)))
}
func F32Sqrt(v float32) float32 {
	return float32(math.Sqrt(float64(v)))
}
func F32Add(a, b float32) float32 {
	return a + b
}
func F32Sub(a, b float32) float32 {
	return a - b
}
func F32Mul(a, b float32) float32 {
	return a * b
}
func F32Div(a, b float32) float32 {
	return a / b
}
func F32Min(a, b float32) float32 {
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		return float32(math.NaN())
	}
	if a == 0 && b == 0 {
		if math.Signbit(float64(a)) || math.Signbit(float64(b)) {
			return float32(math.Copysign(0, -1))
		}
	}
	if a < b {
		return a
	}
	return b
}
func F32Max(a, b float32) float32 {
	if math.IsNaN(float64(a)) || math.IsNaN(float64(b)) {
		return float32(math.NaN())
	}
	if a == 0 && b == 0 {
		if !math.Signbit(float64(a)) || !math.Signbit(float64(b)) {
			return 0.0
		}
		return float32(math.Copysign(0, -1))
	}
	if a > b {
		return a
	}
	return b
}
func F32Copysign(a, b float32) float32 {
	return float32(math.Copysign(float64(a), float64(b)))
}
