// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func I32Eqz(a int32) int32 {
	if a == 0 {
		return 1
	}
	return 0
}
func I32Eq(a, b int32) int32 {
	if a == b {
		return 1
	}
	return 0
}
func I32Ne(a, b int32) int32 {
	if a != b {
		return 1
	}
	return 0
}
func I32Lt_s(a, b int32) int32 {
	if a < b {
		return 1
	}
	return 0
}
func I32Lt_u(a, b int32) int32 {
	if uint32(a) < uint32(b) {
		return 1
	}
	return 0
}
func I32Gt_s(a, b int32) int32 {
	if a > b {
		return 1
	}
	return 0
}
func I32Gt_u(a, b int32) int32 {
	if uint32(a) > uint32(b) {
		return 1
	}
	return 0
}
func I32Le_s(a, b int32) int32 {
	if a <= b {
		return 1
	}
	return 0
}
func I32Le_u(a, b int32) int32 {
	if uint32(a) <= uint32(b) {
		return 1
	}
	return 0
}
func I32Ge_s(a, b int32) int32 {
	if a >= b {
		return 1
	}
	return 0
}
func I32Ge_u(a, b int32) int32 {
	if uint32(a) >= uint32(b) {
		return 1
	}
	return 0
}
