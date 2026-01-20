// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func F32Convert_i32_s(v int32) float32 {
	return float32(v)
}
func F32Convert_i32_u(v int32) float32 {
	return float32(uint32(v))
}
func F32Convert_i64_s(v int64) float32 {
	return float32(v)
}
func F32Convert_i64_u(v int64) float32 {
	return float32(uint64(v))
}

func F64Convert_i32_s(v int32) float64 {
	return float64(v)
}
func F64Convert_i32_u(v int32) float64 {
	return float64(uint32(v))
}
func F64Convert_i64_s(v int64) float64 {
	return float64(v)
}
func F64Convert_i64_u(v int64) float64 {
	return float64(uint64(v))
}

func I32Wrap_i64(v int64) int32 {
	return int32(v)
}
func I64Extend_i32_s(v int32) int64 {
	return int64(v)
}
func I64Extend_i32_u(v int32) int64 {
	return int64(uint32(v))
}
func F32Demote_f64(v float64) float32 {
	return float32(v)
}
func F64Promote_f32(v float32) float64 {
	return float64(v)
}
