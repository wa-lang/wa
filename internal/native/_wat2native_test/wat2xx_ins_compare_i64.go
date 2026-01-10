// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

func I64Eqz(a int64) int32 {
	if a == 0 {
		return 1
	}
	return 0
}
func I64Eq(a, b int64) int32 {
	if a == b {
		return 1
	}
	return 0
}
func I64Ne(a, b int64) int32 {
	if a != b {
		return 1
	}
	return 0
}
func I64Lt_s(a, b int64) int32 {
	if a < b {
		return 1
	}
	return 0
}
func I64Lt_u(a, b int64) int32 {
	if uint64(a) < uint64(b) {
		return 1
	}
	return 0
}
func I64Gt_s(a, b int64) int32 {
	if a > b {
		return 1
	}
	return 0
}
func I64Gt_u(a, b int64) int32 {
	if uint64(a) > uint64(b) {
		return 1
	}
	return 0
}
func I64Le_s(a, b int64) int32 {
	if a <= b {
		return 1
	}
	return 0
}
func I64Le_u(a, b int64) int32 {
	if uint64(a) <= uint64(b) {
		return 1
	}
	return 0
}
func I64Ge_s(a, b int64) int32 {
	if a >= b {
		return 1
	}
	return 0
}
func I64Ge_u(a, b int64) int32 {
	if uint64(a) >= uint64(b) {
		return 1
	}
	return 0
}
