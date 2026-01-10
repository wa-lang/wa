// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func I32Reinterpret_f32(v float32) int32 { return 0 }
func I64Reinterpret_f64(v float64) int64 { return 0 }
func F32Reinterpret_i32(v int32) float32 { return 0 }
func F64Reinterpret_i64(v int64) float64 { return 0 }
