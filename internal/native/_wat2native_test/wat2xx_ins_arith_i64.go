// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

import "math/bits"

func I64Clz(v int64) int64 {
	return int64(bits.LeadingZeros64(uint64(v)))
}
func I64Ctz(v int64) int64 {
	return int64(bits.TrailingZeros64(uint64(v)))
}
func I64Popcnt(v int64) int64 {
	return int64(bits.OnesCount64(uint64(v)))
}
func I64Add(a, b int64) int64 {
	return a + b
}
func I64Sub(a, b int64) int64 {
	return a - b
}
func I64Mul(a, b int64) int64 {
	return a * b
}
func I64Div_s(a, b int64) int64 {
	if b == 0 {
		panic("div zero")
	}
	// I64_MIN / -1 会导致溢出
	if a == -9223372036854775808 && b == -1 {
		return a
	}
	return a / b
}
func I64Div_u(a, b int64) int64 {
	if b == 0 {
		panic("div zero")
	}
	return int64(uint64(a) / uint64(b))
}
func I64Rem_s(a, b int64) int64 {
	if b == 0 {
		panic("div zero")
	}
	if a == -9223372036854775808 && b == -1 {
		return 0
	}
	return a % b
}
func I64Rem_u(a, b int64) int64 {
	if b == 0 {
		panic("div zero")
	}
	return int64(uint64(a) % uint64(b))
}
func I64And(a, b int64) int64 {
	return a & b
}
func I64Or(a, b int64) int64 {
	return a | b
}
func I64Xor(a, b int64) int64 {
	return a ^ b
}
func I64Shl(a, b int64) int64 {
	return a << (uint64(b) & 63)
}
func I64Shr_s(a, b int64) int64 {
	return a >> (uint64(b) & 63)
}
func I64Shr_u(a, b int64) int64 {
	return int64(uint64(a) >> (uint64(b) & 63))
}
func I64Rotl(a, b int64) int64 {
	return int64(bits.RotateLeft64(uint64(a), int(b)))
}
func I64Sotr(a, b int64) int64 {
	return int64(bits.RotateLeft64(uint64(a), -int(b)))
}
