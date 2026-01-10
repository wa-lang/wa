// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

import "math"

func F64Abs(v float64) float64 {
	return math.Float64frombits(math.Float64bits(v) &^ (1 << 63))
}
func F64Neg(v float64) float64 {
	return math.Float64frombits(math.Float64bits(v) ^ (1 << 63))
}
func F64Ceil(v float64) float64 {
	return math.Ceil(v)
}
func F64Floor(v float64) float64 {
	return math.Floor(v)
}
func F64Trunc(v float64) float64 {
	return math.Trunc(v)
}
func F64Nearest(v float64) float64 {
	return math.RoundToEven(v)
}
func F64Sqrt(v float64) float64 {
	return math.Sqrt(v)
}
func F64Add(a, b float64) float64 {
	return a + b
}
func F64Sub(a, b float64) float64 {
	return a - b
}
func F64Mul(a, b float64) float64 {
	return a * b
}
func F64Div(a, b float64) float64 {
	return a / b
}
func F64Min(a, b float64) float64 {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.NaN()
	}
	if a == 0 && b == 0 {
		if math.Signbit(a) || math.Signbit(b) {
			return math.Copysign(0, -1)
		}
	}
	if a < b {
		return a
	}
	return b
}
func F64Max(a, b float64) float64 {
	if math.IsNaN(a) || math.IsNaN(b) {
		return math.NaN()
	}
	if a == 0 && b == 0 {
		if !math.Signbit(a) || !math.Signbit(b) {
			return 0.0
		}
		return math.Copysign(0, -1)
	}
	if a > b {
		return a
	}
	return b
}
func F64Copysign(a, b float64) float64 {
	return math.Copysign(a, b)
}
