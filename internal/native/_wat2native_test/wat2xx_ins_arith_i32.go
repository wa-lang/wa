// Copyright (C) 2026 武汉凹语言科技有限公司
// SPDX-License-Identifier: AGPL-3.0-or-later

//go:build !inline_asm

package main

import "math/bits"

func I32Clz(v int32) int32 {
	return int32(bits.LeadingZeros32(uint32(v)))
}
func I32Ctz(v int32) int32 {
	return int32(bits.TrailingZeros32(uint32(v)))
}
func I32Popcnt(v int32) int32 {
	return int32(bits.OnesCount32(uint32(v)))
}
func I32Add_(a, b int32) int32 {
	return a + b
}
func I32Sub(a, b int32) int32 {
	return a - b
}
func I32Mul(a, b int32) int32 {
	return a * b
}
func I32Div_s(a, b int32) int32 {
	if b == 0 {
		panic("div zero")
	}
	// 补码最小值除以 -1 会溢出
	if a == -2147483648 && b == -1 {
		return a
	}
	return a / b
}
func I32Div_u(a, b int32) int32 {
	if b == 0 {
		panic("div zero")
	}
	return int32(uint32(a) / uint32(b))
}
func I32Rem_s(a, b int32) int32 {
	if b == 0 {
		panic("div zero")
	}
	// 处理最小值溢出的余数
	if a == -2147483648 && b == -1 {
		return 0
	}
	return a % b
}
func I32Rem_u(a, b int32) int32 {
	if b == 0 {
		panic("div zero")
	}
	return int32(uint32(a) % uint32(b))
}
func I32And(a, b int32) int32 {
	return a & b
}
func I32Or(a, b int32) int32 {
	return a | b
}
func I32Xor(a, b int32) int32 {
	return a ^ b
}
func I32Shl(a, b int32) int32 {
	return a << (uint32(b) & 31)
}
func I32Shr_s(a, b int32) int32 {
	return a >> (uint32(b) & 31)
}
func I32Shr_u(a, b int32) int32 {
	return int32(uint32(a) >> (uint32(b) & 31))
}
func I32Rotl(v, c int32) int32 {
	return int32(bits.RotateLeft32(uint32(v), int(c)))
}
func I32Rotr(v, c int32) int32 {
	return int32(bits.RotateLeft32(uint32(v), -int(c)))
}
