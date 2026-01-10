// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func LocalGetI32(offset int) int32   { return 0 }
func LocalGetI64(offset int) int64   { return 0 }
func LocalGetF32(offset int) float32 { return 0 }
func LocalGetF64(offset int) float64 { return 0 }

func LocalSetI32(offset int, v int32)   {}
func LocalSetI64(offset int, v int64)   {}
func LocalSetF32(offset int, v float32) {}
func LocalSetF64(offset int, v float64) {}

func LocalTee(offset int) {}
