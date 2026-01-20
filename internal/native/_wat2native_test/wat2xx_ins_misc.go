// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func Drop() {}

func SelectI32(i int32, a, b int32) int32 {
	if i != 0 {
		return a
	}
	return b
}
func SelectI64(i int32, a, b int64) int64 {
	if i != 0 {
		return a
	}
	return b
}
func SelectF32(i int32, a, b float32) float32 {
	if i != 0 {
		return a
	}
	return b
}
func SelectF64(i int32, a, b float64) float64 {
	if i != 0 {
		return a
	}
	return b
}
