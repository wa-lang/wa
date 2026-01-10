// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func Drop() {}

func SelectI32(i int32, a, b int32) int32     { return 0 }
func SelectI64(i int32, a, b int64) int64     { return 0 }
func SelectF32(i int32, a, b float32) float32 { return 0 }
func SelectF64(i int32, a, b float64) float64 { return 0 }

func I32Wrap_i64()     {}
func I64Extend_i32_s() {}
func I64Extend_i32_u() {}
func F32Demote_f64()   {}
func F64Promote_f32()  {}
